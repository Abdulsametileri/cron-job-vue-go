<template>
  <div>
    <div v-if="jobs.length === 0">You have no alarm</div>
    <div v-for="(job,index) in jobs" :key="index">
      <b-card :title="(index + 1) + '.Alarm: ' + job.name">
        <b-form-select disabled id="repeat" v-model="job.repeatType"
                       :options="repeatTypes"></b-form-select>
        {{job.time}}
        <b-img :src="job.imageUrl" fluid alt="The Image"></b-img>
      </b-card>
    </div>
  </div>
</template>

<script>
import Cookies from "js-cookie";

export default {
  name: "AlarmList",
  data() {
    return {
      jobs: [],
    }
  },
  async created() {
    try {
      let token = Cookies.get('token')
      const res = await fetch('/api/list-alarm?token=' + token, {
        method: 'GET'
      })
      let r = await res.json()
      this.jobs = r.data
    } catch (e) {
      console.error(e)
    }
  }
}
</script>

<style scoped>

</style>