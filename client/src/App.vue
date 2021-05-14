<template>
  <div id="app">
    <b-container class="my-5">
      <b-nav pills align="center">
        <b-nav-item exact exact-active-class="active" :to="{name: 'AlarmList'}">Alarmları Listele</b-nav-item>
        <b-nav-item exact exact-active-class="active" :to="{name: 'AlarmCreate'}">Alarm Ekle</b-nav-item>
      </b-nav>
      <hr>
      <router-view/>
    </b-container>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>

<style lang="scss">

</style>
