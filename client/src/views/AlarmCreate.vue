<template>
  <div>
    <b-form v-if="show" @submit.prevent="onSubmit" @reset.prevent="onReset">
      <b-form-group label="Hatırlatılacak resim" label-for="id">
        <b-form-file
            id="file"
            v-model="form.imgFile"
            placeholder="Resim Seç"
            drop-placeholder="Buraya Bırak..."
        ></b-form-file>
      </b-form-group>

      <b-form-group label="İsim" label-for="name">
        <b-form-input id="name" v-model="form.name"></b-form-input>
      </b-form-group>

      <b-form-group :class="marginTop" label="Saat Seç" label-for="timepicker">
        <b-form-timepicker id="timepicker"
                           locale="tr"
                           placeholder="Saat Seç"
                           v-model="form.time"
                           now-button
                           class="mb-2"></b-form-timepicker>
      </b-form-group>

      <b-form-group :class="marginTop" label="Tekrar" label-for="repeat">
        <b-form-select id="repeat" v-model="form.repeatType"
                       :options="repeatTypes" required></b-form-select>
      </b-form-group>

      <b-button type="submit" variant="primary">Planla</b-button>
      <b-button type="reset" variant="danger" class="ml-3">Sıfırla</b-button>
    </b-form>
  </div>
</template>

<script>

const repeatTypes = [
  {value: 0, text: "Pazartesi"},
  {value: 1, text: "Salı"},
  {value: 2, text: "Çarşamba"},
  {value: 3, text: "Perşembe"},
  {value: 4, text: "Cuma"},
  {value: 5, text: "Cumartesi"},
  {value: 6, text: "Pazar"},
  {value: 7, text: "Her Gün"},
]

const defaultFormItem = {
  name: '',
  imgFile: null,
  time: '',
  repeatType: 7,
}

export default {
  name: "AlarmCreate",
  data() {
    return {
      show: true,
      marginTop: 'mt-3',
      repeatTypes,
      form: {
        ...defaultFormItem
      },
    }
  },
  methods: {
    formValidation() {
      let errorMsg = ""
      if (this.form.name === '')
        errorMsg += "İsim boş bırakılamaz."
      if (this.form.imgFile === null)
        errorMsg += "Dosya bölümü boş bırakılamaz. "
      if (this.form.time === '')
        errorMsg += "Saat boş bırakılamaz. "

      if (errorMsg !== "") {
        this.$bvToast.toast(errorMsg, {
          title: 'Validasyon Hatası',
          variant: 'danger',
        })
        return false
      }
      return true
    },
    onSubmit() {
      if (!this.formValidation())
        return

      console.log(this.form)
    },
    onReset() {
      this.form = {
        ...defaultFormItem
      }
      // Trick to reset/clear native browser form validation state
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