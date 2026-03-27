import { useState, useEffect } from 'react'
import { getAllAppointments, getAvailableDoctors, bookAppointment, bookOTAppointment } from '../../api/appointments'
import { useAuth } from '../../context/AuthContext'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import Modal from '../../components/Modal'
import Alert from '../../components/Alert'
import LoadingSpinner from '../../components/LoadingSpinner'

export default function ReceptionistAppointments() {
  const { user } = useAuth()
  const [appointments,      setAppointments]      = useState([])
  const [loading,           setLoading]           = useState(true)
  const [modal,             setModal]             = useState(null)
  const [alert,             setAlert]             = useState({ type: '', message: '' })
  const [saving,            setSaving]            = useState(false)
  const [availableDoctors,  setAvailableDoctors]  = useState([])
  const [step,              setStep]              = useState(1)

  const [slotForm, setSlotForm] = useState({ date: '', time: '' })
  const [bookForm, setBookForm] = useState({
    patient_email:   '',
    doctor_id:       '',
    receptionist_id: user?.user_id || '',
    date:            '',
    time:            '',
  })

  const fetchData = async () => {
    try {
      const res = await getAllAppointments()
      setAppointments(res.data.appointments || [])
    } catch { setAppointments([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  const showAlert = (type, message) => { setAlert({ type, message }); setTimeout(() => setAlert({ type: '', message: '' }), 5000) }

  const openBook = (type) => {
    setModal(type) // 'regular' or 'ot'
    setStep(1)
    setSlotForm({ date: '', time: '' })
    setBookForm({ patient_email: '', doctor_id: '', receptionist_id: user?.user_id || '', date: '', time: '' })
    setAvailableDoctors([])
  }

  const handleGetDoctors = async () => {
    if (!slotForm.date || !slotForm.time) return showAlert('error', 'Date and time are required')
    setSaving(true)
    try {
      const res = await getAvailableDoctors(slotForm)
      setAvailableDoctors(res.data.available_doctors || [])
      setBookForm({ ...bookForm, date: slotForm.date, time: slotForm.time })
      setStep(2)
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Failed to fetch doctors')
    } finally { setSaving(false) }
  }

  const handleBook = async () => {
    if (!bookForm.patient_email) return showAlert('error', 'Patient email is required')
    if (!bookForm.doctor_id)     return showAlert('error', 'Please select a doctor')
    setSaving(true)
    try {
      const payload = { ...bookForm, doctor_id: parseInt(bookForm.doctor_id), receptionist_id: parseInt(bookForm.receptionist_id) }
      if (modal === 'regular') {
        await bookAppointment(payload)
      } else {
        await bookOTAppointment(payload)
      }
      showAlert('success', `${modal === 'ot' ? 'OT ' : ''}Appointment booked successfully`)
      setModal(null)
      fetchData()
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Booking failed — transaction rolled back')
    } finally { setSaving(false) }
  }

  const columns = [
    { key: 'patient_name', label: 'Patient' },
    { key: 'doctor_name',  label: 'Doctor' },
    { key: 'date',         label: 'Date' },
    { key: 'time',         label: 'Time' },
    { key: 'status', label: 'Status', render: (row) => (
      <span className={row.status === 'pending' ? 'badge-warning' : 'badge-success'}>{row.status}</span>
    )},
    { key: 'ot_id', label: 'OT', render: (row) => (
      <span>{row.ot_id && row.ot_id !== 'None' ? `OT #${row.ot_id}` : '—'}</span>
    )},
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Appointments"
        subtitle={`${appointments.length} total appointments`}
        action={
          <div className="flex gap-2">
            <button onClick={() => openBook('regular')} className="btn-primary">+ Book Appointment</button>
            <button onClick={() => openBook('ot')}      className="btn-secondary">+ Book OT</button>
          </div>
        }
      />

      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />
      <Table columns={columns} data={appointments} emptyMessage="No appointments found" />
{modal === 'regular' && (
  <Modal title="Book Appointment" onClose={() => setModal(null)}>
    {step === 1 ? (
      <div className="space-y-4">
        <p className="text-sm text-gray-500">Step 1 — Select date and time to see available doctors</p>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Date</label>
          <input type="date" className="input-field" value={slotForm.date} onChange={e => setSlotForm({...slotForm, date: e.target.value})} />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Time</label>
          <input type="time" className="input-field" value={slotForm.time} onChange={e => setSlotForm({...slotForm, time: e.target.value})} />
        </div>
        <button onClick={handleGetDoctors} disabled={saving} className="btn-primary w-full">
          {saving ? 'Checking...' : 'Check Available Doctors →'}
        </button>
      </div>
    ) : (
      <div className="space-y-4">
        <p className="text-sm text-gray-500">Step 2 — Select a doctor and enter patient details</p>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Patient Email</label>
          <input className="input-field" value={bookForm.patient_email} onChange={e => setBookForm({...bookForm, patient_email: e.target.value})} placeholder="patient@email.com" />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Available Doctors — {slotForm.date} at {slotForm.time}
          </label>
          {availableDoctors.length === 0 ? (
            <p className="text-red-500 text-sm">No doctors available at this slot</p>
          ) : (
            <div className="space-y-2 max-h-48 overflow-y-auto">
              {availableDoctors.map(doc => (
                <label key={doc.id} className={`flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors ${bookForm.doctor_id == doc.id ? 'border-primary-500 bg-primary-50' : 'border-gray-200 hover:bg-gray-50'}`}>
                  <input
                    type="radio"
                    name="doctor"
                    value={doc.id}
                    checked={bookForm.doctor_id == doc.id}
                    onChange={e => setBookForm({...bookForm, doctor_id: e.target.value})}
                  />
                  <div>
                    <p className="font-medium text-sm">{doc.name}</p>
                    <p className="text-xs text-gray-500">{doc.specialization} • {doc.email}</p>
                  </div>
                </label>
              ))}
            </div>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Receptionist ID</label>
          <input type="number" className="input-field" value={bookForm.receptionist_id} onChange={e => setBookForm({...bookForm, receptionist_id: e.target.value})} />
        </div>
        <div className="flex gap-3 pt-2">
          <button onClick={() => setStep(1)} className="btn-secondary flex-1">← Back</button>
          <button onClick={handleBook} disabled={saving || availableDoctors.length === 0} className="btn-primary flex-1">
            {saving ? 'Booking...' : 'Confirm Booking'}
          </button>
        </div>
      </div>
    )}
  </Modal>
)}

{modal === 'ot' && (
  <Modal title="Book OT Appointment" onClose={() => setModal(null)}>
    <div className="space-y-4">
      <div className="p-3 bg-blue-50 rounded-lg text-sm text-blue-700 border border-blue-200">
        ℹ️ Book OT as requested by the doctor. The system will automatically assign an available OT room.
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Patient Email</label>
        <input
          className="input-field"
          value={bookForm.patient_email}
          onChange={e => setBookForm({...bookForm, patient_email: e.target.value})}
          placeholder="patient@email.com"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Doctor ID</label>
        <input
          type="number"
          className="input-field"
          value={bookForm.doctor_id}
          onChange={e => setBookForm({...bookForm, doctor_id: e.target.value})}
          placeholder="Doctor who requested the OT"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Receptionist ID</label>
        <input
          type="number"
          className="input-field"
          value={bookForm.receptionist_id}
          onChange={e => setBookForm({...bookForm, receptionist_id: e.target.value})}
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Date</label>
        <input
          type="date"
          className="input-field"
          value={bookForm.date}
          onChange={e => setBookForm({...bookForm, date: e.target.value})}
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Time</label>
        <input
          type="time"
          className="input-field"
          value={bookForm.time}
          onChange={e => setBookForm({...bookForm, time: e.target.value})}
        />
      </div>
      <div className="flex gap-3 pt-2">
        <button
          onClick={handleBook}
          disabled={saving}
          className="btn-primary flex-1"
        >
          {saving ? 'Booking...' : 'Book OT Appointment'}
        </button>
        <button onClick={() => setModal(null)} className="btn-secondary flex-1">Cancel</button>
      </div>
    </div>
  </Modal>
)}
      
    </div>
  )
}
