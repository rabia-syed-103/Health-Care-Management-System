import { useState, useEffect } from 'react'
import { getMyAppointments } from '../../api/doctor'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import StatCard from '../../components/StatCard'
import LoadingSpinner from '../../components/LoadingSpinner'
import { useAuth } from '../../context/AuthContext'

export default function DoctorDashboard() {
  const { user } = useAuth()
  const [appointments, setAppointments] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetch = async () => {
      try {
        const res = await getMyAppointments()
        setAppointments(res.data.appointments || [])
      } catch { setAppointments([]) }
      finally { setLoading(false) }
    }
    fetch()
  }, [])

  const pending   = appointments.filter(a => a.status === 'pending').length
  const completed = appointments.filter(a => a.status === 'completed').length
  const withOT    = appointments.filter(a => a.ot && a.ot !== 'None').length

  const columns = [
    { key: 'patient_name', label: 'Patient' },
    { key: 'blood_group',  label: 'Blood Group' },
    { key: 'p_no',         label: 'Phone' },
    { key: 'date',         label: 'Date' },
    { key: 'time',         label: 'Time' },
    { key: 'status', label: 'Status', render: (row) => (
      <span className={row.status === 'pending' ? 'badge-warning' : 'badge-success'}>
        {row.status}
      </span>
    )},
    { key: 'ot', label: 'OT', render: (row) => (
      <span>{row.ot && row.ot !== 'None' ? `OT #${row.ot}` : '—'}</span>
    )},
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title={`Welcome, ${user?.name || 'Doctor'}`}
        subtitle="Your appointments and schedule"
      />

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
