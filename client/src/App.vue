<template>
  <div id="app">
    <b-container class="my-5">
      <template v-if="isAuthenticated">
        <b-nav pills align="center">
          <b-nav-item exact exact-active-class="active" :to="{name: 'AlarmList'}">{{ $t("listAlarm") }}</b-nav-item>
          <b-nav-item exact exact-active-class="active" :to="{name: 'AlarmCreate'}">{{ $t("createAlarm") }}</b-nav-item>
        </b-nav>
        <hr>
        <router-view/>
      </template>
      <template v-else>
        <h5>{{ $t('telegramBotMsg') }}</h5>
        <a target="_blank" href="https://t.me/reminder_1996_bot">Bot Link</a>

        <div class="mt-3"></div>

        <b-form-group label="Token" label-for="token">
          <b-form-input id="token" v-model="token"></b-form-input>
        </b-form-group>

        <b-button @click="addToken" variant="outline-primary">OK</b-button>
      </template>
    </b-container>
  </div>
</template>

<script>
import Cookies from 'js-cookie'

export default {
  name: 'app',
  data() {
    return {
      isAuthenticated: false,
      token: '',
    }
  },
  methods: {
    async addToken() {
      if (this.token === ""){
        alert('Invalid token')
        return
      }

      try {
        const { status } = await fetch("/api/validate-token?token=" + this.token, {
          method: "GET",
        })
        if (status === 200) {
          this.isAuthenticated = true
          Cookies.set("authenticate", true, {expires: 365})
        } else {
          alert('Invalid token')
        }
      } catch (e) {
        console.error(e)
        alert(e.message)
      }
    }
  },
  created() {
    this.isAuthenticated = Cookies.get("authenticate")
    console.log(Cookies.get("authenticate"))
  }
}
</script>

<style lang="scss">

</style>
