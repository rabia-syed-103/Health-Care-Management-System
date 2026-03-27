import { useState, useEffect } from 'react'
import { getAllAppointments } from '../../api/appointments'
import { getAllPatients } from '../../api/patients'
import { getAllERPatients } from '../../api/er'
import PageHeader from '../../components/PageHeader'
import StatCard from '../../components/StatCard'
import LoadingSpinner from '../../components/LoadingSpinner'
import { useAuth } from '../../context/AuthContext'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

export default function ReceptionistDashboard() {
  const { user } = useAuth()
  const [stats,   setStats]   = useState({ appointments: 0, patients: 0, er: 0, pending: 0 })
  const [chartData, setChartData] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchAll = async () => {
      try {
        const [appts, patients, er] = await Promise.all([
          getAllAppointments(),
          getAllPatients(),
          getAllERPatients(),
        ])
        const appointments = appts.data.appointments || []
        const pending = appointments.filter(a => a.status === 'pending').length

        setStats({
          appointments: appointments.length,
          patients:     patients.data.patients?.length || 0,
          er:           er.data.er_patients?.length    || 0,
          pending,
        })

        // Group appointments by date for chart
        const byDate = {}
        appointments.forEach(a => {
          byDate[a.date] = (byDate[a.date] || 0) + 1
        })
        const chartArr = Object.entries(byDate)
          .sort(([a], [b]) => a.localeCompare(b))
          .slice(-7)
          .map(([date, count]) => ({ date: date.slice(5), count }))
        setChartData(chartArr)
      } catch { }
      finally { setLoading(false) }
    }
    fetchAll()
  }, [])

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title={`Welcome, ${user?.email || 'Receptionist'}`}
        badge={user?.id ? `ID: ${user.id}` : null}
        subtitle="Hospital front desk overview"
      />

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <StatCard title="Total Appointments" value={stats.appointments} color="blue"   />
        <StatCard title="Pending"            value={stats.pending}      color="yellow" />
        <StatCard title="Total Patients"     value={stats.patients}     color="green"  />
        <StatCard title="ER Patients"        value={stats.er}           color="red"    />
      </div>

      {chartData.length > 0 && (
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Appointments (Last 7 Days)</h2>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" tick={{ fontSize: 12 }} />
              <YAxis allowDecimals={false} />
              <Tooltip />
              <Bar dataKey="count" fill="#14b8a6" radius={[4,4,0,0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  )
}
