import api from './axios'

export const getMyAppointments  = ()      => api.get('/doctor/my-appointments')
export const getPatientHistory  = (email) => api.get(`/doctor/patient-history/${email}`)