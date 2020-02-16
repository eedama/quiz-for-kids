import axios from 'axios'


/*
export default axios.create({
    baseURL: `https://opentdb.com`,
    withCredentials: false,
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    }
})
*/

export default axios.create({
    baseURL: `https://api.eedama.org`,
    withCredentials: false,
    headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    }
})