<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
import {
  mapGetters,
  mapActions
} from 'vuex'
import requestInterval from 'request-interval'

export default {
  name: 'App',
  computed: {
    ...mapGetters(['nodes'])
  },
  mounted() {
    this['blockchain/getLastBlock']()
    this['blockchain/subNewBlock']()
    this['blockchain/subRoundStep']()
    requestInterval(1000, () => this['blockchain/getConsensusState']())
    this['blockchain/getStatus']()
    this['blockchain/getNodes']()
    this['blockchain/getValidators']()
  },
  methods: {
    ...mapActions([
      'blockchain/getLastBlock',
      'blockchain/subNewBlock',
      'blockchain/subRoundStep',
      'blockchain/getConsensusState',
      'blockchain/getStatus',
      'blockchain/getNodes',
      'blockchain/getValidators'
    ])
  }
}
</script>
