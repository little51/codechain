<template>
  <div class="account-container">
    <p class="account-propmt">{{ prompt }}</p>
    <el-button
      class="account-el-button"
      type="primary"
      @click="addNewAccount()"
    >Add New Account</el-button>
    <textarea
      v-if="temp_new_account.length > 0"
      v-model="json_data"
      class="account-input"
      readonly
    />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  beforeRouteLeave(to, from, next) {
    if (this.json_data === '') {
      next()
    } else {
      if (global.confirm('please recording the information, you can\'t find them after leaving page')) {
        this['account/delete_new_account']()
        next()
      }
    }
  },
  data() {
    return {
      temp_new_account: [],
      json_data: ``
    }
  },
  computed: {
    ...mapGetters([
      'account'
    ]),
    prompt() {
      return this.temp_new_account.length > 0
        ? 'you had created a new account, now you can jump the page of Assets to Register assets'
        : 'we found you have not any account, please click the button of "Add New Account" to add a new account  temporarily, don\'t refresh web page during this period'
    }
  },
  mounted() {
    // 如果store中有数据的话，就展示，并且不能让
    if (this.account.address !== '') {
      this.objToArray()
      this.el_button_disabled = true
      this.json_data = this.jsonStructure()
    }
  },
  methods: {
    ...mapActions(['account/add_new_account', 'account/delete_new_account']),

    // 将从store中拿到的account转化到temp_new_account的数组当中展示
    objToArray() {
      this.temp_new_account.push({ Account_key: 'address', Account_value: this.account.address })
      // this.temp_new_account.push({ Account_key: 'error', Account_value: this.account.error })
      // this.temp_new_account.push({ Account_key: 'privateKey', Account_value: this.account.privateKey })
      // this.temp_new_account.push({ Account_key: 'publicKey', Account_value: this.account.publicKey })
    },
    async addNewAccount() {
      await this['account/add_new_account']()
      this.objToArray()
      this.json_data = this.jsonStructure()
    },
    // 构造json数据格式
    jsonStructure() {
      return `
        {\n
            address: "${this.account.address}",\n
            error: "${this.account.error}",\n
            privateKey: "${this.account.privateKey}",\n
            publicKey: "${this.account.publicKey}"\n
        }
      `
    }
  }
}
</script>

<style lang="scss" scoped>
.account {
  &-propmt {
    font-size: 15px;
    color: #97A8BE;
  }
  &-container {
    margin-left: 10px;
    margin-right: 30px;
    margin-top: 20px;
  }
  &-el-button {
    margin-bottom: 10px;
  }
  &-input {
    width: 100%;
    height: 250px;
    border-radius: 15px;
    background-color: #f4f5f7;
    font-weight: 500;
    color: #304156;
    font-family:'Lucida Sans', 'Lucida Sans Regular', 'Lucida Grande', 'Lucida Sans Unicode', Geneva, Verdana, sans-serif;
    border-style: none;
    box-shadow: 1px 1px 5px 1px #888888
  }
}
</style>
