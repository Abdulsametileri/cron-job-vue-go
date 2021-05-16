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
                           :locale="appLocal"
                           :placeholder="$t('timePickerPlaceHolder')"
                           v-model="form.time"
                           now-button
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
  name: 'asd',
  imgFile: null,
  time: '22:22',
  repeatType: 3,
}

export default {
  name: "AlarmCreate",
  data() {
    return {
      show: true,
      marginTop: 'mt-3',
      repeatTypes: [
        {value: -1, text: this.$t('weekDays.default')},
        {value: 1, text: this.$t('weekDays.monday')},
        {value: 2, text: this.$t('weekDays.tuesday')},
        {value: 3, text: this.$t('weekDays.wednesday')},
        {value: 4, text: this.$t('weekDays.thursday')},
        {value: 5, text: this.$t('weekDays.friday')},
        {value: 6, text: this.$t('weekDays.saturday')},
        {value: 0, text: this.$t('weekDays.sunday')},
        {value: 7, text: this.$t('weekDays.all')},
      ],
      form: {
        ...defaultFormItem
      },
    }
  },
  methods: {
    formValidation() {
      if (this.isDevelopment)
        return true

      let errorMsg = ""
      if (this.form.name === '')
        errorMsg += this.$t('formError.name') + " "
      if (this.form.imgFile === null)
        errorMsg += this.$t('formError.imgFile') + " "
      if (this.form.time === '')
        errorMsg += this.$t('formError.time') + " "
      if (this.form.repeatType === -1)
        errorMsg += this.$t('formError.repeatType')

      if (errorMsg !== "") {
        this.$bvToast.toast(errorMsg, {
          title: 'Error',
          variant: 'danger',
        })
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
        let k = await res.json()
        console.log(k)

        if (res.status === 200) {
          this.$router.replace('/')
          return
        }

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