import axios from 'axios'

export default axios.create({
    baseURL: `https://quizz.eedama.org`,
    withCredentials: false,
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    }
})