import api from './axios'

export const getAllERPatients  = ()     => api.get('/receptionist/er-patients')
export const addERPatient      = (data) => api.post('/receptionist/er-patients', data)
export const getAllEREntries    = ()     => api.get('/receptionist/er-entries')
export const addEREntry        = (data) => api.post('/receptionist/er-entries', data)