import Vue from 'vue'
import Vuex from 'vuex'

// import getters from './getters'
import * as getters from './getters'
// import * as actions from './actions'

import app from './modules/app'
import settings from './modules/settings'
import user from './modules/user'
import blockchain from './modules/blockchain'
import config from './modules/config'
import account from './modules/account'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    app,
    settings,
    user,
    blockchain,
    config,
    account
  },
  getters
  // actions
})

export default store
