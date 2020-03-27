<template>
  <div class="account-container">
    <el-button type="primary">Asset signature</el-button>
    <el-button type="primary">Add New Account</el-button>
    <el-table
      :data="account"
      style="width: 100%;margin-bottom: 20px;"
      :header-cell-style="{background:'#ECEEF1'}"
      row-key="id"
    >
      <el-table-column
        prop="TestnetDate_key"
        label="Testnet Data"
        width="180"
        show-header="false"
      />
      <el-table-column
        prop="TestnetDate_value"
      />
    </el-table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {

}
</script>

<style lang="scss" scoped>
.account {
  &-container {
    margin-left: 10px;
    margin-right: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
