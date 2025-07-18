<script>
import { SECRET } from '@shell/config/types';
import { LabeledInput } from '@components/Form/LabeledInput';
import LabeledSelect from '@shell/components/form/LabeledSelect';
import UnitInput from '@shell/components/form/UnitInput';

export default {
  components: {
    LabeledInput,
    LabeledSelect,
    UnitInput,
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
        processorSockets: this.value.processorSockets ? parseInt(this.value.processorSockets) : "",
        processorCores: this.value.processorCores ? parseInt(this.value.processorCores) : "",
        memory: this.value.memory ? parseInt(this.value.memory) : "",
        memoryBalloon: this.value.memoryBalloon ? parseInt(this.value.memoryBalloon) : "",
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
    isoDeviceSelectOptions(){
      if(this.isoDeviceSelectOptions == null) {
        return
      }

      if(this.isoDeviceSelectOptions.length == 1) {
        this.currentValue.isoDevice = this.isoDeviceSelectOptions[0];
      }

      if(!this.isoDeviceSelectOptions.includes(this.currentValue.isoDevice)) {
        this.currentValue.isoDevice = '';
      }
    },
    networkInterfaceSelectOptions(){
      if(this.networkInterfaceSelectOptions == null) {
        return
      }

      if(this.networkInterfaceSelectOptions.length == 1) {
        this.currentValue.networkInterface = this.networkInterfaceSelectOptions[0];
      }

      if(!this.networkInterfaceSelectOptions.includes(this.currentValue.networkInterface)) {
        this.currentValue.networkInterface = '';
      }
    },
  },
  methods: {
    validate(){
      // Default optional number fields to empty string
      if(this.currentValue.processorSockets === null || this.currentValue.processorSockets === 0) {
        this.currentValue.processorSockets = "";
      }

      if(this.currentValue.processorCores === null || this.currentValue.processorCores === 0) {
        this.currentValue.processorCores = "";
      }

      if(this.currentValue.memory === null || this.currentValue.memory === 0) {
        this.currentValue.memory = "";
      }

      if(this.currentValue.memoryBalloon === null) {
        this.currentValue.memoryBalloon = "";
      }

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

      if(this.currentValue.processorSockets != "" && this.currentValue.processorSockets < 1) {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.processorCores != "" && this.currentValue.processorCores < 1) {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.memory != "" && this.currentValue.memory < 1) {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.memoryBalloon != "" && this.currentValue.memoryBalloon < 0) {
        this.$emit('validationChanged', false);
        return
      }

      if(this.currentValue.memory != "" && this.currentValue.memoryBalloon != "" && this.currentValue.memory < this.currentValue.memoryBalloon) {
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
      this.value.processorSockets = this.currentValue.processorSockets.toString();
      this.value.processorCores = this.currentValue.processorCores.toString();
      this.value.memory = this.currentValue.memory.toString();
      this.value.memoryBalloon = this.currentValue.memoryBalloon.toString();

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
      if(this.credential == null) {
        this.templates = null;
        return;
      }

      if(this.resourcePools && this.currentValue.resourcePool == '') {
        this.templates = [];
        return;
      }

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
      if(this.credential == null) {
        this.devices = null;
        return;
      }

      if(this.templates && this.currentValue.template == '') {
        this.devices = [];
        return;
      }

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
        
        this.devices = data;
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
    isoDeviceSelectOptions() {
      if(this.devices == null) {
        return null;
      }

      return this.devices
        .filter(device => device.key)
        .filter(device => (device.value && typeof device.value == "string"))
        .filter(device => device.value.includes('media=cdrom'))
        .map(device => device.key);
    },
    networkInterfaceSelectOptions() {
      if(this.devices == null) {
        return null;
      }

      return this.devices
        .filter(device => device.key)
        .filter(device => device.key.startsWith('net'))
        .map(device => device.key);
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
      :value="currentValue.resourcePool"
      @change="e => { currentValue.resourcePool = e.target.value}"
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
      v-if="templates == null"
      class="mt-20"
      type="number"
      :mode="mode"
      :disabled="disabled"
      :value="currentValue.template"
      @change="e => { currentValue.template = e.target.value}"
      label-key="cluster.machineConfig.pve.template.label"
      required
      min="1"
      step="1"
    />

    <LabeledSelect
      v-else
      class="mt-20"
      :mode="mode"
      :disabled="disabled || (resourcePools != null && !currentValue.resourcePool)"
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
          v-if="devices == null"
          type="text"
          :mode="mode"
          :disabled="disabled"
          :value="currentValue.isoDevice"
          @change="e => { currentValue.isoDevice = e.target.value}"
          label-key="cluster.machineConfig.pve.devices.iso.label"
          tooltip-key="cluster.machineConfig.pve.devices.iso.tooltip"
          required
        />

        <LabeledSelect
          v-else
          :mode="mode"
          :disabled="disabled || (templates != null && !currentValue.template)"
          v-model:value="currentValue.isoDevice"
          :options="isoDeviceSelectOptions"
          label-key="cluster.machineConfig.pve.devices.iso.label"
          tooltip-key="cluster.machineConfig.pve.devices.iso.tooltip"
          required
        />
      </div>
      <div class="col span-6">
        <!-- Network interface -->
        <LabeledInput
          v-if="devices == null"
          type="text"
          :mode="mode"
          :value="currentValue.networkInterface"
          @change="e => { currentValue.networkInterface = e.target.value}"
          label-key="cluster.machineConfig.pve.devices.network.label"
          tooltip-key="cluster.machineConfig.pve.devices.network.tooltip"
          required
        />

        <LabeledSelect
          v-else
          :mode="mode"
          :disabled="disabled || (templates != null && !currentValue.template)"
          v-model:value="currentValue.networkInterface"
          :options="networkInterfaceSelectOptions"
          label-key="cluster.machineConfig.pve.devices.network.label"
          tooltip-key="cluster.machineConfig.pve.devices.network.tooltip"
          required
        />
      </div>
    </div>

    <h2 class="mt-20">
      <t k="cluster.machineConfig.pve.hardware.header" />
    </h2>
    <div class="row mt-20">
      <div class="col span-6">
        <!-- Processor sockets -->
        <UnitInput
          type="number"
          :mode="mode"
          v-model:value="currentValue.processorSockets"
          label-key="cluster.machineConfig.pve.hardware.processorSockets.label"
          suffix="sockets"
          min="0"
          step="1"
        />
      </div>
      <div class="col span-6">
        <!-- Processor cores -->
        <UnitInput
          type="number"
          :mode="mode"
          v-model:value="currentValue.processorCores"
          label-key="cluster.machineConfig.pve.hardware.processorCores.label"
          suffix="cores"
          min="0"
          step="1"
        />
      </div>
    </div>

    <div class="row mt-20">
      <div class="col span-6">
        <!-- Memory -->
        <UnitInput
          type="number"
          :mode="mode"
          v-model:value="currentValue.memory"
          label-key="cluster.machineConfig.pve.hardware.memory.label"
          suffix="MiB"
          :status="currentValue.memory != '' && currentValue.memoryBalloon != '' && currentValue.memoryBalloon > currentValue.memory ? 'error' : undefined"
          min="0"
          step="256"
        />
      </div>
      <div class="col span-6">
        <!-- Memory balloon -->
        <UnitInput
          type="number"
          :mode="mode"
          v-model:value="currentValue.memoryBalloon"
          label-key="cluster.machineConfig.pve.hardware.memoryBalloon.label"
          suffix="MiB"
          :status="currentValue.memory != '' && currentValue.memoryBalloon != '' && currentValue.memoryBalloon > currentValue.memory ? 'error' : undefined"
          min="0"
          step="256"
          :max="currentValue.memory != '' ? currentValue.memory : undefined"
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
