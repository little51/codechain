<template>
  <div class="account-container">
    <p class="account-propmt">{{ prompt }}</p>
    <el-button
      class="account-el-button"
      :type="el_button_account_type"
      :disabled="el_button_disabled"
      @click="addNewAccount()"
    >Add New Account</el-button>
    <span v-if="temp_new_account.length > 0"> &gt; </span>
    <el-button
      v-if="temp_new_account.length > 0"
      type="primary"
    >
      <router-link :to="{path:`/Assets/new/`}">
        Add New Assets
      </router-link>
    </el-button>
    <el-table
      v-if="temp_new_account.length > 0"
      :data="temp_new_account"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="Account_key"
        label="Account Prop"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="Account_value"
      />
    </el-table>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  data() {
    return {
      temp_new_account: [],
      el_button_account_type: 'primary',
      el_button_disabled: false
    }
  },
  computed: {
    ...mapGetters([
      'account'
    ]),
    prompt() {
      return this.temp_new_account.length > 0
        ? 'you had created a new account, now you can click the button of "Add New Assets" to jump the page of Assets'
        : 'we found you have not any account, please click the button of "Add New Account" to add a new account  temporarily, don\'t refresh web page during this period'
    }
  },
  mounted() {
    // 如果store中有数据的话，就展示，并且不能让
    if (this.account.address !== '') {
      this.objToArray()
      this.el_button_account_type = 'info'
      this.el_button_disabled = true
    }
  },
  methods: {
    ...mapActions(['account/add_new_account']),

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
      this.el_button_account_type = 'info'
      this.el_button_disabled = true
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
}
</style>
