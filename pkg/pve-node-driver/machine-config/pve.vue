<script>
import { LabeledInput } from '@components/Form/LabeledInput';

export default {
  components: {
    LabeledInput,
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
  },
  data(){
    let data = {
      resourcePool: this.value.resourcePool ?? '',
      template: parseInt(this.value.template ?? '0'),
      isoDevice: this.value.isoDevice ?? '',
      networkInterface: this.value.networkInterface ?? '',
      sshUser: this.value.sshUser ?? '',
      sshPort: parseInt(this.value.sshPort ?? '0'),
    }

    if(data.sshUser == ''){
      data.sshUser = 'service'
    }

    if(data.sshPort == 0){
      data.sshPort = 22
    }

    return data
  },
  created(){
    this.updateValue();
  },
  methods: {
    updateValue(){
      if(this.resourcePool == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.template < 1) {
        this.$emit('validationChanged', false);
        return
      }

      if(this.isoDevice == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.networkInterface == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.sshUser == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.sshPort < 1) {
        this.$emit('validationChanged', false);
        return
      }

      this.value.resourcePool = this.resourcePool;
      this.value.template = this.template.toString();
      this.value.isoDevice = this.isoDevice;
      this.value.networkInterface = this.networkInterface;
      this.value.sshUser = this.sshUser;
      this.value.sshPort = this.sshPort.toString();

      this.$emit('validationChanged', true);
    },
    test(){
      this.updateValue();
    },
  }
}
</script>

<template>
  <div>
    <div class="mt-20 mb-20">
      <LabeledInput
        type="text"
        :mode="mode"
        v-model:value="resourcePool"
        @change="updateValue"
        label-key="cluster.machineConfig.pve.resourcePool.label"
        required
      />
    </div>

    <div class="mb-20">
      <LabeledInput
        type="number"
        :mode="mode"
        v-model:value="template"
        @change="updateValue"
        label-key="cluster.machineConfig.pve.template.label"
        required
        min="0"
        step="1"
      />
    </div>

    <div class="mb-20">
      <LabeledInput
        type="text"
        :mode="mode"
        v-model:value="isoDevice"
        @change="updateValue"
        label-key="cluster.machineConfig.pve.isoDevice.label"
        required
      />
    </div>

    <div class="mb-20">
      <LabeledInput
        type="text"
        :mode="mode"
        v-model:value="networkInterface"
        @change="updateValue"
        label-key="cluster.machineConfig.pve.networkInterface.label"
        required
      />
    </div>

    <div class="mb-20">
      <LabeledInput
        type="text"
        :mode="mode"
        v-model:value="sshUser"
        @change="updateValue"
        label-key="cluster.machineConfig.pve.sshUser.label"
        required
      />
    </div>

    <div class="mb-20">
      <LabeledInput
        type="number"
        :mode="mode"
        v-model:value="sshPort"
        @change="updateValue"
        label-key="cluster.machineConfig.pve.sshPort.label"
        required
        min="0"
        step="1"
      />
    </div>
  </div>
</template>
