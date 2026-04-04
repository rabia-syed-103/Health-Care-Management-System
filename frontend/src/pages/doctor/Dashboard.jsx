import Alert from '../../components/Alert'
import { useState, useEffect } from 'react'
import { getMyAppointments } from '../../api/doctor'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import StatCard from '../../components/StatCard'
import LoadingSpinner from '../../components/LoadingSpinner'
import { useAuth } from '../../context/AuthContext'
import { cancelAppointment } from '../../api/appointments'
import { exportPDF } from '../../utils/pdfExport'

export default function DoctorDashboard() {
  const { user } = useAuth()
  const [appointments, setAppointments] = useState([])
  const [loading, setLoading] = useState(true)
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [cancelling, setCancelling] = useState(null)

  useEffect(() => {
    const fetch = async () => {
      try {
        const res = await getMyAppointments()
        setAppointments(res.data.appointments || [])
        console.log('DATA:', res.data.appointments)
      } catch { setAppointments([]) }
      finally { setLoading(false) }
    }
    fetch()
  }, [])
const pending   = appointments.filter(a => a.status === 'pending').length
const completed = appointments.filter(a => a.status === 'completed').length
const withOT    = appointments.filter(a => a.ot_id && a.ot_id !== 'None' && a.status === 'pending').length
  
  const handleCancel = async (appointmentId) => {
  console.log('Cancelling appointment ID:', appointmentId)
  setCancelling(appointmentId)
  try {
    await cancelAppointment(appointmentId)
    setAlert({ type: 'success', message: 'Appointment cancelled successfully' })
    const res = await getMyAppointments()
    setAppointments(res.data.appointments || [])
  } catch (err) {
    setAlert({ type: 'error', message: err.response?.data?.error || 'Failed to cancel appointment' })
  } finally { setCancelling(null) }
}
const downloadPDF = () => exportPDF(
  'My Appointments Report',
  ['Patient', 'Blood Group', 'Phone', 'Email', 'Date', 'Time', 'Status', 'OT'],
  appointments.map(r => [
    r.patient_name, r.patient_blood_group, r.patient_phone,
    r.patient_email, r.date, r.time, r.status,
    r.ot_id && r.ot_id !== 'None' ? `OT #${r.ot_id}` : '—'
  ]),
  'my-appointments-report'
)
const columns = [
  { key: 'patient_name', label: 'Patient' },
  { key: 'patient_blood_group',  label: 'Blood Group' },
  { key: 'patient_phone',         label: 'Phone' },
  {key: 'patient_email',         label: 'Email' },
  { key: 'date',         label: 'Date' },
  { key: 'time',         label: 'Time' },
  { key: 'status', label: 'Status', render: (row) => (
    <span className={
      row.status === 'pending'   ? 'badge-warning' :
      row.status === 'cancelled' ? 'badge-danger'  : 'badge-success'
    }>
      {row.status}
    </span>
  )},
  { key: 'ot_id', label: 'OT', render: (row) => (
    <span>{row.ot_id && row.ot_id !== 'None' ? `OT #${row.ot_id}` : '—'}</span>
  )},
  { key: 'actions', label: 'Actions', render: (row) => (
    row.status === 'pending' ? (
      <button
        onClick={() => handleCancel(row.id)}
disabled={cancelling === row.id}
        className="text-red-600 hover:underline text-sm disabled:opacity-50"
      >
        {cancelling === row.id ? 'Cancelling...' : 'Cancel'}

      </button>
    ) : <span className="text-gray-400 text-sm">—</span>
  )},
]

  if (loading) return <LoadingSpinner />
  console.log('USER OBJECT:', user)
  console.log('USER ID:', user?.id)
  

  return (
    <div>
  <PageHeader
    title={`Welcome, ${user?.email || 'Doctor'}`}
    subtitle="Your appointments and schedule"
    badge={user?.id ? `ID: ${user.id}` : null}
    action={<button onClick={downloadPDF} className="btn-secondary">⬇ Download PDF</button>}
  />
    <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

      <div className="grid grid-cols-3 gap-4 mb-8">
        <StatCard title="Total Appointments" value={appointments.length} color="blue"  />
        <StatCard title="Pending"            value={pending}            color="yellow" />
        <StatCard title="With OT"            value={withOT}             color="purple" />
      </div>

      <div className="card">
        <h2 className="text-lg font-semibold text-gray-800 mb-4">My Appointments</h2>
        <Table columns={columns} data={appointments} emptyMessage="No appointments found" />
      </div>
    </div>
  )
}


