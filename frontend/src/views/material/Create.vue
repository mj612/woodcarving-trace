<template>
  <div class="material-create">
    <div class="page-container">
      <div class="card">
        <h2 class="form-title">添加原料</h2>
        
        <el-form :model="form" :rules="rules" ref="form" label-width="120px">
          <el-form-item label="木材种类" prop="woodType">
            <el-select v-model="form.woodType" placeholder="请选择木材种类">
              <el-option label="红木" value="红木"></el-option>
              <el-option label="黄花梨" value="黄花梨"></el-option>
              <el-option label="紫檀" value="紫檀"></el-option>
              <el-option label="花梨木" value="花梨木"></el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="产地" prop="origin">
            <el-input v-model="form.origin" placeholder="请输入产地"></el-input>
          </el-form-item>
          
          <el-form-item label="采伐证编号" prop="harvestCert">
            <el-input v-model="form.harvestCert" placeholder="请输入采伐证编号"></el-input>
          </el-form-item>
          
          <el-form-item label="数量(立方米)" prop="quantity">
            <el-input-number v-model="form.quantity" :min="0" :step="0.1"></el-input-number>
          </el-form-item>
          
          <el-form-item label="质量等级" prop="quality">
            <el-select v-model="form.quality" placeholder="请选择质量等级">
              <el-option label="优" value="优"></el-option>
              <el-option label="良" value="良"></el-option>
              <el-option label="一般" value="一般"></el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="描述">
            <el-input 
              v-model="form.description" 
              type="textarea" 
              :rows="4"
              placeholder="请输入原料描述"
            ></el-input>
          </el-form-item>
          
          <el-form-item label="证书文件">
            <el-upload
              class="upload-demo"
              action="#"
              :auto-upload="false"
              :on-change="handleFileChange"
            >
              <el-button size="small" type="primary">点击上传</el-button>
              <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
            </el-upload>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="submitForm" :loading="loading">提交</el-button>
            <el-button @click="resetForm">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { createMaterial, uploadFile } from '@/api'

export default {
  name: 'MaterialCreate',
  data() {
    return {
      loading: false,
      form: {
        woodType: '',
        origin: '',
        harvestCert: '',
        quantity: 0,
        quality: '',
        description: '',
        images: '',
        certFile: ''
      },
      rules: {
        woodType: [
          { required: true, message: '请选择木材种类', trigger: 'change' }
        ],
        origin: [
          { required: true, message: '请输入产地', trigger: 'blur' }
        ],
        harvestCert: [
          { required: true, message: '请输入采伐证编号', trigger: 'blur' }
        ],
        quantity: [
          { required: true, message: '请输入数量', trigger: 'blur' }
        ],
        quality: [
          { required: true, message: '请选择质量等级', trigger: 'change' }
        ]
      },
      uploading: false
    }
  },
  methods: {
    submitForm() {
      this.$refs.form.validate(async (valid) => {
        if (valid) {
          this.loading = true
          try {
            await createMaterial(this.form)
            this.$message.success('原料添加成功')
            this.$router.push('/dashboard/materials')
          } catch (error) {
            console.error('添加原料失败:', error)
            this.$message.error('添加原料失败')
          } finally {
            this.loading = false
          }
        } else {
          this.$message.error('请填写完整信息')
          return false
        }
      })
    },
    
    resetForm() {
      this.$refs.form.resetFields()
    },
    
    async handleFileChange(file) {
      // 上传文件到服务器
      this.uploading = true
      try {
        const res = await uploadFile(file.raw)
        if (res.code === 200) {
          this.form.certFile = res.data.url
          this.$message.success('文件上传成功')
        }
      } catch (error) {
        console.error('文件上传失败:', error)
        this.$message.error('文件上传失败')
        this.form.certFile = ''
      } finally {
        this.uploading = false
      }
    }
  }
}
</script>

<style scoped>
.upload-demo {
  width: 360px;
}
</style>