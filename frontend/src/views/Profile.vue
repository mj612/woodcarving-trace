<template>
  <div class="profile">
    <div class="page-container">
      <div class="card">
        <h2 class="form-title">个人资料</h2>
        
        <el-form :model="form" label-width="100px">
          <el-form-item label="用户名">
            <el-input v-model="form.username" disabled></el-input>
          </el-form-item>
          
          <el-form-item label="真实姓名">
            <el-input v-model="form.realName"></el-input>
          </el-form-item>
          
          <el-form-item label="手机号">
            <el-input v-model="form.phone"></el-input>
          </el-form-item>
          
          <el-form-item label="邮箱">
            <el-input v-model="form.email"></el-input>
          </el-form-item>
          
          <el-form-item label="公司">
            <el-input v-model="form.company"></el-input>
          </el-form-item>
          
          <el-form-item label="地址">
            <el-input v-model="form.address" type="textarea"></el-input>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="saveProfile">保存修改</el-button>
            <el-button @click="resetForm">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 修改密码 -->
      <div class="card">
        <h2 class="form-title">修改密码</h2>
        
        <el-form :model="passwordForm" label-width="100px">
          <el-form-item label="原密码" prop="oldPassword">
            <el-input v-model="passwordForm.oldPassword" type="password"></el-input>
          </el-form-item>
          
          <el-form-item label="新密码" prop="newPassword">
            <el-input v-model="passwordForm.newPassword" type="password"></el-input>
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input v-model="passwordForm.confirmPassword" type="password"></el-input>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="changePassword">修改密码</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { updateUserInfo, changePassword } from '@/api'

export default {
  name: 'Profile',
  data() {
    return {
      form: {
        username: '',
        realName: '',
        phone: '',
        email: '',
        company: '',
        address: ''
      },
      passwordForm: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    }
  },
  computed: {
    ...mapState(['userInfo'])
  },
  mounted() {
    // 确保用户信息已加载
    if (!this.userInfo) {
      this.$store.dispatch('getUserInfo').then(() => {
        this.loadUserInfo()
      }).catch(() => {
        this.$message.error('获取用户信息失败')
        this.$router.push('/login')
      })
    } else {
      this.loadUserInfo()
    }
  },
  methods: {
    loadUserInfo() {
      if (!this.userInfo) {
        return
      }
      
      this.form = {
        username: this.userInfo.username || '',
        realName: this.userInfo.realName || '',
        phone: this.userInfo.phone || '',
        email: this.userInfo.email || '',
        company: this.userInfo.company || '',
        address: this.userInfo.address || ''
      }
    },
    
    async saveProfile() {
      try {
        await updateUserInfo({
          realName: this.form.realName,
          phone: this.form.phone,
          email: this.form.email,
          company: this.form.company,
          address: this.form.address
        })
        
        // 更新本地存储的用户信息
        this.$store.commit('SET_USER_INFO', {
          ...this.userInfo,
          ...this.form
        })
        
        this.$message.success('个人信息保存成功')
      } catch (error) {
        console.error('保存个人信息失败:', error)
        this.$message.error('保存失败，请重试')
      }
    },
    
    async changePassword() {
      if (this.passwordForm.newPassword !== this.passwordForm.confirmPassword) {
        this.$message.error('两次输入的密码不一致')
        return
      }
      
      if (this.passwordForm.newPassword.length < 6) {
        this.$message.error('新密码长度不能少于6位')
        return
      }
      
      if (!this.passwordForm.oldPassword) {
        this.$message.error('请输入原密码')
        return
      }
      
      try {
        await changePassword({
          oldPassword: this.passwordForm.oldPassword,
          newPassword: this.passwordForm.newPassword
        })
        
        this.$message.success('密码修改成功，请重新登录')
        
        // 清空表单
        this.passwordForm = {
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        }
        
        // 2秒后跳转到登录页
        setTimeout(() => {
          // 清除token和用户信息
          this.$store.dispatch('logout')
          this.$router.push('/login')
        }, 2000)
        
      } catch (error) {
        console.error('修改密码失败:', error)
        // 根据错误信息显示不同的提示
        if (error.response && error.response.data && error.response.data.msg) {
          this.$message.error(error.response.data.msg)
        } else {
          this.$message.error('修改失败，请检查原密码是否正确')
        }
      }
    },
    
    resetForm() {
      this.loadUserInfo()
    }
  }
}
</script>

<style scoped>
.profile ::v-deep .el-form-item__label {
  font-weight: bold;
}
</style>