<template>
  <div>
    <b-form v-if="show" @submit.prevent="onSubmit" @reset.prevent="onReset">
      <b-form-group :label="$t('image')" label-for="id">
        <b-form-file
            id="file"
            v-model="form.imgFile"
            :placeholder="$t('fileUploadPlaceholder')"
            :drop-placeholder="$t('fileDropPlaceholder')"
        ></b-form-file>
      </b-form-group>

      <b-form-group :label="$t('documentName')" label-for="name">
        <b-form-input id="name" v-model="form.name"></b-form-input>
      </b-form-group>

      <b-form-group :class="marginTop"
                    :label="$t('time')" label-for="timepicker">
        <b-form-timepicker id="timepicker"
                           :locale="$i18n.locale"
                           :placeholder="$t('timePickerPlaceHolder')"
                           v-model="form.time"
                           class="mb-2"></b-form-timepicker>
      </b-form-group>

      <b-form-group :class="marginTop" :label="$t('repeatOptions')" label-for="repeat">
        <b-form-select id="repeat" v-model="form.repeatType"
                       :options="repeatTypes"></b-form-select>
      </b-form-group>

      <b-button type="submit" variant="primary">{{ $t('createAlarmBtn') }}</b-button>
      <b-button type="reset" variant="danger" class="ml-3">{{ $t('resetAlarmForm') }}</b-button>
    </b-form>
  </div>
</template>

<script>
import Cookies from 'js-cookie'

const defaultFormItem = {
  name: '',
  imgFile: null,
  time: '',
  repeatType: -1,
}

export default {
  name: "AlarmCreate",
  data() {
    return {
      show: true,
      marginTop: 'mt-3',
      form: {
        ...defaultFormItem
      },
    }
  },
  methods: {
    formValidation() {
      let errorMsg = ""
      if (this.form.name === '')
        errorMsg += this.$t('formError.name') + "\n"
      if (this.form.imgFile === null)
        errorMsg += this.$t('formError.imgFile') + "\n"
      if (this.form.time === '')
        errorMsg += this.$t('formError.time') + "\n"
      if (this.form.repeatType === -1)
        errorMsg += this.$t('formError.repeatType')

      if (errorMsg !== "") {
        this.showErrorMessage(errorMsg)
        return false
      }
      return true
    },
    async onSubmit() {
      if (!this.formValidation())
        return

      try {
        const formData = new FormData()
        formData.append("token", Cookies.get('token'))
        formData.append("name", this.form.name)
        formData.append("file", this.form.imgFile)
        formData.append("fileType", this.form.imgFile.type)
        formData.append("fileName", this.form.imgFile.name)
        formData.append("time", this.form.time)
        formData.append("repeatType", this.form.repeatType)

        const res = await fetch("/api/create-alarm", {
          method: "POST",
          body: formData
        })
        let {code, message} = await res.json()

        if (code === 200) {
          this.$router.replace('/')
          return
        }
        this.showErrorMessage(message)
      } catch (e) {
        console.error(e)
      }
    },
    onReset() {
      this.form = {
        ...defaultFormItem
      }
      this.show = false
      this.$nextTick(() => {
        this.show = true
      })
    }
  }
}
</script>

<style scoped>

</style>