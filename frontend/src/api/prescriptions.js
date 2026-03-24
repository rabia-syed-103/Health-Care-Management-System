import api from './axios'

export const prescribeMedicines      = (data) => api.post('/prescriptions/prescribe', data)
export const dispenseMedicines       = (data) => api.post('/dispensing/dispense', data)
export const getPendingPrescriptions = ()     => api.get('/pharmacist/pending-prescriptions')