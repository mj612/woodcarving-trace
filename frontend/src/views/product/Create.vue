<template>
  <div class="product-create">
    <div class="page-container">
      <div class="card">
        <h2 class="form-title">创建产品</h2>
        
        <el-form :model="form" :rules="rules" ref="form" label-width="120px">
          <el-form-item label="产品名称" prop="productName">
            <el-input v-model="form.productName" placeholder="请输入产品名称"></el-input>
          </el-form-item>
          
          <el-form-item label="选择原料" prop="materialId">
            <el-select v-model="form.materialId" placeholder="请选择原料" filterable style="width: 100%;">
              <el-option 
                v-for="material in availableMaterials" 
                :key="material.materialId"
                :label="`${material.materialId} - ${material.woodType} (${material.origin})`"
                :value="material.materialId">
                <span>{{ material.materialId }}</span>
                <span style="margin-left: 10px; color: #8492a6;">{{ material.woodType }} - {{ material.origin }}</span>
              </el-option>
            </el-select>
            <div v-if="availableMaterials.length === 0" style="color: #999; font-size: 12px; margin-top: 5px;">
              暂无可用原料，请先到"原料管理"页面添加原料
            </div>
          </el-form-item>
          
          <el-form-item label="尺寸(长x宽x高)" prop="dimensions">
            <el-input v-model="form.dimensions" placeholder="例如：50x30x20cm"></el-input>
          </el-form-item>
          
          <el-form-item label="重量(kg)" prop="weight">
            <el-input-number v-model="form.weight" :min="0" :step="0.1"></el-input-number>
          </el-form-item>
          
          <el-form-item label="雕刻工艺描述">
            <el-input 
              v-model="form.craftDesc" 
              type="textarea" 
              :rows="4"
              placeholder="请输入雕刻工艺描述"
            ></el-input>
          </el-form-item>
          
          <el-form-item label="雕刻耗时(小时)">
            <el-input-number v-model="form.carveTime" :min="0"></el-input-number>
          </el-form-item>
          
          <el-form-item label="产品图片">
            <el-upload
              class="upload-demo"
              action="#"
              :auto-upload="false"
              :on-change="handleImageChange"
              list-type="picture"
            >
              <el-button size="small" type="primary">点击上传</el-button>
              <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过2MB</div>
            </el-upload>
          </el-form-item>
          
          <el-form-item label="设计图纸">
            <el-upload
              class="upload-demo"
              action="#"
              :auto-upload="false"
              :on-change="handleDesignChange"
            >
              <el-button size="small" type="primary">点击上传</el-button>
              <div slot="tip" class="el-upload__tip">支持PDF、DWG等格式</div>
            </el-upload>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="submitForm" :loading="loading">创建产品</el-button>
            <el-button @click="resetForm">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { createProduct, getMaterialList, uploadFile } from '@/api'

export default {
  name: 'ProductCreate',
  data() {
    return {
      loading: false,
      form: {
        productName: '',
        materialId: '',
        dimensions: '',
        weight: 0,
        craftDesc: '',
        carveTime: 0,
        images: [],
        designFile: null
      },
      availableMaterials: [],
      rules: {
        productName: [
          { required: true, message: '请输入产品名称', trigger: 'blur' }
        ],
        materialId: [
          { required: true, message: '请选择原料', trigger: 'change' }
        ],
        dimensions: [
          { required: true, message: '请输入产品尺寸', trigger: 'blur' }
        ],
        weight: [
          { required: true, message: '请输入产品重量', trigger: 'blur' }
        ]
      }
    }
  },
  
  mounted() {
    this.loadMaterials()
  },
  
  methods: {
    async loadMaterials() {
      try {
        const res = await getMaterialList()
        console.log('原料列表响应:', res)
        
        // API返回的是 { list: [], total: 0 } 格式
        const list = (res.data && res.data.list) ? res.data.list : []
        console.log('原料列表数据:', list)
        
        if (list.length === 0) {
          this.$message.warning('暂无可用原料，请先添加原料')
          this.availableMaterials = []
          return
        }
        
        // 从链上数据中提取字段，显示所有原料（不过滤状态）
        this.availableMaterials = list.map(item => {
          const chainData = item.chainData || {}
          return {
            materialId: item.materialId || chainData.materialId,
            woodType: chainData.woodType || '未知',
            origin: chainData.origin || '',
            quantity: chainData.quantity || 0,
            quality: chainData.quality || '',
            status: chainData.status || 'unknown'
          }
        })
        
        console.log('可用原料列表:', this.availableMaterials)
        
        if (this.availableMaterials.length === 0) {
          this.$message.warning('暂无可用原料')
        }
      } catch (error) {
        console.error('加载原料列表失败:', error)
        this.$message.error('加载原料列表失败: ' + (error.message || '未知错误'))
        this.availableMaterials = []
      }
    },
    
    submitForm() {
      this.$refs.form.validate(async (valid) => {
        if (valid) {
          this.loading = true
          try {
            await createProduct(this.form)
            this.$message.success('产品创建成功')
            this.$router.push('/dashboard/products')
          } catch (error) {
            console.error('创建产品失败:', error)
            this.$message.error('创建产品失败')
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
    
    async handleImageChange(file, fileList) {
      // 上传图片到服务器
      try {
        const res = await uploadFile(file.raw)
        if (res.code === 200) {
          // 收集所有已上传图片的URL
          const uploadedUrls = []
          for (const f of fileList) {
            if (f.raw) {
              // 如果是新上传的文件，上传到服务器
              try {
                const uploadRes = await uploadFile(f.raw)
                if (uploadRes.code === 200) {
                  uploadedUrls.push(uploadRes.data.url)
                }
              } catch (err) {
                console.error('上传失败:', err)
              }
            }
          }
          this.form.images = uploadedUrls.join(',')
          this.$message.success('图片上传成功')
        }
      } catch (error) {
        console.error('图片上传失败:', error)
        this.$message.error('图片上传失败')
      }
    },
    
    async handleDesignChange(file) {
      // 上传设计文件到服务器
      try {
        const res = await uploadFile(file.raw)
        if (res.code === 200) {
          this.form.designFile = res.data.url
          this.$message.success('设计文件上传成功')
        }
      } catch (error) {
        console.error('设计文件上传失败:', error)
        this.$message.error('设计文件上传失败')
        this.form.designFile = ''
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