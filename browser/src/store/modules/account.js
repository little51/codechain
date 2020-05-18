import axios from 'axios'

const defaultAccount = {
  address: '',
  error: '',
  privateKey: '',
  publicKey: ''
}

const state = {
  tem_new_account: defaultAccount
}

const actions = {
  async add_new_account({ commit }) {
    let json = await axios.post(`http://localhost:4000/account/new`, {})
    let result = json.data
    commit('set_tem_new_account', result)
    return Promise.resolve()
  },
  delete_new_account({ commit }) {
    commit('delete_tem_new_accunt')
  }
}

const mutations = {
  set_tem_new_account(state, account) {
    console.log(`正在创建新的账户`)
    state.tem_new_account = account
  },
  delete_tem_new_accunt(state) {
    console.log(`正在清除临时账户`)
    state.tem_new_account = defaultAccount
  }
}

export default {
  namespaced: true,
  state,
  actions,
  mutations
}
