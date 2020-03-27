import axios from 'axios'

const state = {
  tem_new_account: {
    address: '',
    error: '',
    privateKey: '',
    publicKey: ''
  }
}

const actions = {
  async add_new_account({ commit }) {
    let json = await axios.post(`http://localhost:3000/account/new`, {})
    let result = json.data
    commit('set_tem_new_account', result)
    return Promise.resolve()
  }
}

const mutations = {
  set_tem_new_account(state, account) {
    console.log(`正在创建新的账户`)
    state.tem_new_account = account
  }
}

export default {
  namespaced: true,
  state,
  actions,
  mutations
}
