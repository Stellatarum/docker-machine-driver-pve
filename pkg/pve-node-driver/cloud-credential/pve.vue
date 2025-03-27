<script>
import Banner from '@components/Banner/Banner';
import { Checkbox } from '@components/Form/Checkbox';
import { LabeledInput } from '@components/Form/LabeledInput';

export default {
  components: {
    Banner,
    Checkbox,
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
  created() {
    this.value.setData('url', this.value.decodedData.url ?? '');
    this.value.setData('insecureTls', this.value.decodedData.insecureTls ?? false);
    this.value.setData('tokenId', this.value.decodedData.tokenId ?? '');
    this.value.setData('tokenSecret', this.value.decodedData.tokenSecret ?? '');

    this.validate();
  },
  methods: {
    validate(){
      if(this.value.decodedData.url == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.value.decodedData.tokenId == '') {
        this.$emit('validationChanged', false);
        return
      }

      if(this.value.decodedData.tokenSecret == '') {
        this.$emit('validationChanged', false);
        return
      }

      this.$emit('validationChanged', true);
    },
  }
}
</script>

<template>
  <div>
    <div class="mb-20">
      <LabeledInput
        type="text"
        :mode="mode"
        :value="value.decodedData.url"
        @change="e => {value.setData('url', e.target.value); validate();}"
        label-key="cluster.credential.pve.url.label"
        placeholder-key="cluster.credential.pve.url.placeholder"
        required
      />
    </div>
    <div class="mb-20">
      <Checkbox
        :mode="mode"
        :value="value.decodedData.insecureTls"
        @click="e => {value.setData('insecureTls', !value.decodedData.insecureTls); validate();}"
        label-key="cluster.credential.pve.insecureTLS.label"
      />
      <Banner
        v-if="value.decodedData.insecureTls"
        color="warning"
        label-key="cluster.credential.pve.insecureTLS.warning"
      />
    </div>
    <div class="mb-20">
      <LabeledInput
        type="text"
        :mode="mode"
        :value="value.decodedData.tokenId"
        @change="e => {value.setData('tokenId', e.target.value); validate();}"
        label-key="cluster.credential.pve.tokenID.label"
        placeholder-key="cluster.credential.pve.tokenID.placeholder"
        tooltip-key="cluster.credential.pve.tokenID.tooltip"
        required
      />
    </div>
    <div>
      <LabeledInput
        type="password"
        :mode="mode"
        :value="value.decodedData.tokenSecret"
        @change="e => {value.setData('tokenSecret', e.target.value); validate();}"
        label-key="cluster.credential.pve.tokenSecret.label"
        placeholder-key="cluster.credential.pve.tokenSecret.placeholder"
        required
      />
    </div>
  </div>
</template>
