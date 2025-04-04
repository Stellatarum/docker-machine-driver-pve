<script>
import { SECRET } from '@shell/config/types';
import { LabeledInput } from '@components/Form/LabeledInput';
import LabeledSelect from '@shell/components/form/LabeledSelect';

export default {
  components: {
    LabeledInput,
    LabeledSelect,
  },
  props: {
    mode: {
      type:     String,
      required: true,
    },
    value: {
      type:     Object,
      required: true,
    },
    credentialId: {
      type:     String,
      required: true,
    },
  },
  data(){
    return {
      // Count of currently running external data fetches.
      fetchingCount: 0,

      // Cloud credential data.
      credential: null,

      // Available Proxmox VE resource pools.
      resourcePools: null,

      // Available Proxmox VE templates in current resource pool.
      templates: null,

      // Available devices on current Proxmox VE template.
      devices: null,

      // Current input values.
      currentValue: {
        resourcePool: this.value.resourcePool ?? '',
        template: this.value.template ? parseInt(this.value.template) : 0,
        isoDevice: this.value.isoDevice ?? '',
        networkInterface: this.value.networkInterface ?? '',
        sshUser: this.value.sshUser ? this.value.sshUser : 'service',
        sshPort: this.value.sshPort ? parseInt(this.value.sshPort) : 22,
      },
    }
  },
  created(){
    this.validate();
  },
  async fetch(){
    await this.fetchCredential();
    await this.fetchResourcePools();
    await this.fetchTemplates();
    await this.fetchDevices();
  },
  watch: {
    async credentialId(){
      await this.$fetch();
      this.validate();
    },
    'currentValue.resourcePool': async function (){
      await this.fetchTemplates();
      await this.fetchDevices();
    },
    'currentValue.template': async function (){
      await this.fetchDevices();
    },
    currentValue: {
      deep: true,
      handler(){
        this.validate();
      },
    },
    resourcePools(){
      if(this.resourcePools == null) {
        return
      }

      if(this.resourcePools.length == 1) {
        this.currentValue.resourcePool = this.resourcePools[0];
      }

      if(!this.resourcePools.includes(this.currentValue.resourcePool)) {
        this.currentValue.resourcePool = '';
      }
    },
    templates(){
      if(this.templates == null) {
        return
      }

      if(this.templates.length == 1) {
        this.currentValue.template = this.templates[0].vmid;
      }

      const currentTemplate = this.templates.find(template => this.currentValue.template == template.vmid)

      if(!currentTemplate) {
        this.currentValue.template = 0;
      }
    },
    devices(){
      if(this.devices == null) {
        return
      }

      if(this.devices.length == 1) {
        this.currentValue.isoDevice = this.devices[0];
      }

      if(!this.devices.includes(this.currentValue.isoDevice)) {
        this.currentValue.isoDevice = '';
      }
    },
  },
  methods: {
    validate(){
      // Sanity check current value
      if(this.currentValue.resourcePool == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.template < 1) {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.isoDevice == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.networkInterface == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.sshUser == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.sshPort < 1) {
        this.$emit('validationChanged', false);
        return
      }

      // Copy current value to the 'value' prop
      this.value.resourcePool = this.currentValue.resourcePool;
      this.value.template = this.currentValue.template.toString();
      this.value.isoDevice = this.currentValue.isoDevice;
      this.value.networkInterface = this.currentValue.networkInterface;
      this.value.sshUser = this.currentValue.sshUser;
      this.value.sshPort = this.currentValue.sshPort.toString();

      this.$emit('validationChanged', true);
    },
    async fetchCredential(){
      try {
        this.fetchingCount += 1;

        const secret = await this.$store.dispatch('management/find', {
          type: SECRET,
          id: this.credentialId.replace(':', '/'),
        });

        const decodedData = {
          url: atob(secret.data['pvecredentialConfig-url']),
          insecureTls: atob(secret.data['pvecredentialConfig-insecureTls']) == 'true',
          tokenId: atob(secret.data['pvecredentialConfig-tokenId']),
          tokenSecret: atob(secret.data['pvecredentialConfig-tokenSecret']),
        }

        if(decodedData.insecureTls) {
          // If insecure TLS is enabled then we don't have ability to call Proxmox VE from the UI.
          throw new Error('Insecure TLS is enabled')
        }
        
        this.credential = decodedData;
      } catch(e) {
        console.warn("Failed to fetch Proxmox VE credential, enhanced form fields are disabled", e)
        this.credential = null;
      } finally {
        this.fetchingCount -= 1;
      }
    },
    async fetchResourcePools(){
      try {
        this.fetchingCount += 1;

        const { data } = await this.fetchFromProxmox('/api2/json/pools');

        this.resourcePools = data.map(pool => pool.poolid);
      } catch(e) {
        this.resourcePools = null;
      } finally {
        this.fetchingCount -= 1;
      }
    },
    async fetchTemplates(){
      try {
        this.fetchingCount += 1;

        if(!this.currentValue.resourcePool) {
          throw new Error('Resource pool is not selected')
        }

        const { data } = await this.fetchFromProxmox(`/api2/json/pools?type=qemu&poolid=${this.currentValue.resourcePool}`);
        
        this.templates = data[0].members.filter(vm => vm.template == 1);
      } catch(e) {
        this.templates = null;
      } finally {
        this.fetchingCount -= 1;
      }
    },
    async fetchDevices(){
      try {
        this.fetchingCount += 1;

        if(!this.currentValue.template) {
          throw new Error('Template is not selected');
        }

        const template = this.templates.find(template => this.currentValue.template == template.vmid);

        if(!template) {
          throw new Error("Template not found")
        }

        const { data } = await this.fetchFromProxmox(`/api2/json/nodes/${template.node}/qemu/${template.vmid}/pending`);
        
        this.devices = data
          .filter(device => device.key)
          .filter(device => (device.value && typeof device.value == "string"))
          .filter(device => device.value.includes('media=cdrom'))
          .map(device => device.key);
      } catch(e) {
        this.devices = null;
      } finally {
        this.fetchingCount -= 1;
      }
    },
    fetchFromProxmox(apiPath){
      if(this.credential == null) {
        throw new Error("Credential not available")
      }

      return this.$store.dispatch('management/request', {
        method: 'GET',
        url: '/meta/proxy/' + this.credential.url.replace(/^https?:\/\//, '') + apiPath,
        headers: {
          "X-API-Auth-Header": `PVEAPIToken=${this.credential.tokenId}=${this.credential.tokenSecret}`,
        },
        redirectUnauthorized: false,
      })
    },
    selectTemplate(option) {
      if(Object.keys(this.templateSelectOptions).length == 0) {
        this.currentValue.template = 0;
        return;
      }

      for (const [vmid, name] of Object.entries(this.templateSelectOptions)) {
        if(option == name) {
          this.currentValue.template = vmid;
          return;
        }
      }

      this.currentValue.template = 0;
    },
  },
  computed: {
    disabled(){
      return this.fetchingCount > 0
    },
    templateSelectOptions() {
      if(this.templates === null) {
        return {};
      }

      return this.templates.reduce(
        (templates, template) => ({...templates, [template.vmid]: `${template.name} (${template.vmid})`}),
        {},
      );
    },
    templateSelectValue() {
      return this.templateSelectOptions[this.currentValue.template] ?? '';
    },
  },
}
</script>

<template>
  <div class="mt-20 mb-20">
    <!-- Resource pool  -->
    <LabeledInput
      v-if="resourcePools == null"
      type="text"
      :mode="mode"
      :disabled="disabled"
      v-model:value="currentValue.resourcePool"
      label-key="cluster.machineConfig.pve.resourcePool.label"
      required
    />

    <LabeledSelect
      v-else
      :mode="mode"
      :disabled="disabled"
      v-model:value="currentValue.resourcePool"
      :options="resourcePools"
      label-key="cluster.machineConfig.pve.resourcePool.label"
      required
    />

    <!-- Template -->
    <LabeledInput
      v-if="resourcePools == null && templates == null"
      class="mt-20"
      type="number"
      :mode="mode"
      :disabled="disabled"
      v-model:value="currentValue.template"
      label-key="cluster.machineConfig.pve.template.label"
      required
      min="1"
      step="1"
    />
      
    <LabeledSelect
      v-else
      class="mt-20"
      :mode="mode"
      :disabled="disabled || !currentValue.resourcePool"
      v-model:value="templateSelectValue"
      @option:selected="selectTemplate"
      :options="Object.values(templateSelectOptions)"
      label-key="cluster.machineConfig.pve.template.label"
      required
    />

    <h2 class="mt-20">
      <t k="cluster.machineConfig.pve.devices.header" />
    </h2>
    <div class="row mt-20">
      <div class="col span-6">
        <!-- ISO Device -->
        <LabeledInput
          v-if="resourcePools == null && templates == null && devices == null"
          type="text"
          :mode="mode"
          :disabled="disabled"
          v-model:value="currentValue.isoDevice"
          label-key="cluster.machineConfig.pve.devices.iso.label"
          tooltip-key="cluster.machineConfig.pve.devices.iso.tooltip"
          required
        />
        
        <LabeledSelect
          v-else
          :mode="mode"
          :disabled="disabled || !currentValue.resourcePool || !currentValue.template"
          v-model:value="currentValue.isoDevice"
          :options="devices"
          label-key="cluster.machineConfig.pve.devices.iso.label"
          tooltip-key="cluster.machineConfig.pve.devices.iso.tooltip"
          required
        />
      </div>
      <div class="col span-6">
        <!-- Network interface -->
        <LabeledInput
          type="text"
          :mode="mode"
          v-model:value="currentValue.networkInterface"
          label-key="cluster.machineConfig.pve.devices.network.label"
          tooltip-key="cluster.machineConfig.pve.devices.network.tooltip"
          required
        />
      </div>
    </div>

    <h2 class="mt-20">
      <t k="cluster.machineConfig.pve.ssh.header" />
    </h2>
    <div class="row mt-20">
      <div class="col span-6">
        <!-- SSH Username -->
        <LabeledInput
          type="text"
          :mode="mode"
          v-model:value="currentValue.sshUser"
          label-key="cluster.machineConfig.pve.ssh.username.label"
          tooltip-key="cluster.machineConfig.pve.ssh.username.tooltip"
          required
        />
      </div>
      <div class="col span-6">
        <!-- SSH Port -->
        <LabeledInput
          type="number"
          :mode="mode"
          v-model:value="currentValue.sshPort"
          label-key="cluster.machineConfig.pve.ssh.port.label"
          tooltip-key="cluster.machineConfig.pve.ssh.port.tooltip"
          required
          min="1"
          step="1"
        />
      </div>
    </div>

  </div>
</template>
