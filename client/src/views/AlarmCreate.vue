<template>
  <div>
    <b-form v-if="show" @submit.prevent="onSubmit" @reset.prevent="onReset">
      <h5>{{ $t('telegramBotMsg') }}</h5>
      <a target="_blank" href="https://t.me/reminder_1996_bot">Bot Link</a>

      <div class="mt-3"></div>

      <b-form-group label="Token" label-for="token">
        <b-form-input id="token" v-model="form.token"></b-form-input>
      </b-form-group>

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
                       :options="repeatTypes" required></b-form-select>
      </b-form-group>

      <b-button type="submit" variant="primary">{{ $t('createAlarmBtn') }}</b-button>
      <b-button type="reset" variant="danger" class="ml-3">{{ $t('resetAlarmForm') }}</b-button>
    </b-form>
  </div>
</template>

<script>
const defaultFormItem = {
  token: '',
  name: '',
  imgFile: null,
  time: '00:00:00',
  repeatType: null,
}

export default {
  name: "AlarmCreate",
  data() {
    return {
      show: true,
      marginTop: 'mt-3',
      repeatTypes: [
        {value: null, text: this.$t('weekDays.default')},
        {value: 0, text: this.$t('weekDays.sunday')},
        {value: 1, text: this.$t('weekDays.monday')},
        {value: 2, text: this.$t('weekDays.tuesday')},
        {value: 3, text: this.$t('weekDays.wednesday')},
        {value: 4, text: this.$t('weekDays.thursday')},
        {value: 5, text: this.$t('weekDays.friday')},
        {value: 6, text: this.$t('weekDays.saturday')},
        {value: 7, text: this.$t('weekDays.all')},
      ],
      form: {
        ...defaultFormItem
      },
    }
  },
  methods: {
    formValidation() {
      let errorMsg = ""
      if (this.form.name === '')
        errorMsg += this.$t('formError.name') + " "
      if (this.form.imgFile === null)
        errorMsg += this.$t('formError.imgFile') + " "
      if (this.form.time === '')
        errorMsg += this.$t('formError.time') + " "
      if (this.form.repeatType === null)
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
      if (!this.isDevelopment && !this.formValidation())
        return

      const formData = new FormData()
      formData.append("name", this.form.name)
      formData.append("file", this.form.imgFile)
      formData.append("time", this.form.time)
      formData.append("repeatType", this.form.repeatType)

      await fetch("/api/create-alarm", {
        method: "POST",
        body: formData
      })
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