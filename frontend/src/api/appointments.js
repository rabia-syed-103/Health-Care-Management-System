import api from './axios'

export const getAvailableDoctors  = (data) => api.post('/appointments/available-doctors', data)
export const bookAppointment      = (data) => api.post('/appointments/book', data)
export const bookOTAppointment    = (data) => api.post('/appointments/book-ot', data)
export const getAllAppointments    = ()     => api.get('/receptionist/appointments')