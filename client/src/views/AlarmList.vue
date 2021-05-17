<template>
  <div>
    <div v-if="jobs.length === 0">You have no alarm</div>
    <div v-for="(job,index) in jobs" :key="index">
      <b-card>
        <b-card-title>
          Alarm: {{ job.name }}
        </b-card-title>

        <b-card-sub-title class="float-right">
          <b-button variant="danger" @click="deleteAlarm(job.tag)">
            DELETE
            <b-icon icon="trash"></b-icon>
          </b-button>
        </b-card-sub-title>

        <b-card-text>
          <b-icon icon="alarm" variant="success"/>
          {{ indexStrToWeekDay[job.repeatType] }} ~/~ {{ job.time }}
          <b-icon icon="alarm" variant="success"/>
        </b-card-text>

        <hr>

        <div class="my-4"></div>

        <b-img :src="job.imageUrl"></b-img>
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
  methods: {
    async deleteAlarm(jobTag) {
      try {
        const res = await fetch("/api/delete-alarm?tag=" + jobTag, {
          method: 'POST'
        })
        const {code, message} = await res.json()
        if (code === 200) {
          await this.getAllValidAlarms()
          this.showSucessMessage(this.$t('operationSuccess'))
        } else {
          this.showErrorMessage(message)
        }
      } catch (e) {
        console.error(e)
      }
    },
    async getAllValidAlarms() {
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
  },
  async created() {
    await this.getAllValidAlarms()
  }
}
</script>

<style scoped>

</style>