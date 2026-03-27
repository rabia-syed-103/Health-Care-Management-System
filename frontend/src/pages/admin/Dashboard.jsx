import { useState, useEffect } from 'react'
import { getAllDoctors, getAllReceptionists, getAllPharmacists, getAllBloodManagers } from '../../api/admin'
import { getAllPatients } from '../../api/patients'
import { getAllDonors } from '../../api/donors'
import PageHeader from '../../components/PageHeader'
import StatCard from '../../components/StatCard'
import LoadingSpinner from '../../components/LoadingSpinner'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from 'recharts'

export default function AdminDashboard() {
  const [stats, setStats] = useState({
    doctors: 0, receptionists: 0, pharmacists: 0,
    blood_managers: 0, patients: 0, donors: 0
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchAll = async () => {
      try {
        const [d, r, p, bm, pat, don] = await Promise.all([
          getAllDoctors(),
          getAllReceptionists(),
          getAllPharmacists(),
          getAllBloodManagers(),
          getAllPatients(),
          getAllDonors(),
        ])
        setStats({
          doctors:       d.data.doctors?.length       || 0,
          receptionists: r.data.receptionists?.length || 0,
          pharmacists:   p.data.pharmacists?.length   || 0,
          blood_managers:bm.data.blood_managers?.length || 0,
          patients:      pat.data.patients?.length    || 0,
          donors:        don.data.donors?.length      || 0,
        })
      } catch (err) {
        console.error(err)
      } finally {
        setLoading(false)
      }
    }
    fetchAll()
  }, [])

  const staffData = [
    { name: 'Doctors',       value: stats.doctors },
    { name: 'Receptionists', value: stats.receptionists },
    { name: 'Pharmacists',   value: stats.pharmacists },
    { name: 'Blood Managers',value: stats.blood_managers },
  ]

  const pieData = [
    { name: 'Patients', value: stats.patients },
    { name: 'Donors',   value: stats.donors },
  ]

  const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6']

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader title="Admin Dashboard" subtitle="Hospital overview and statistics" 
      badge={user?.id ? `ID: ${user.id}` : null}/>

      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 mb-8">
        <StatCard title="Doctors"        value={stats.doctors}        color="blue"   />
        <StatCard title="Receptionists"  value={stats.receptionists}  color="teal"   />
        <StatCard title="Pharmacists"    value={stats.pharmacists}    color="green"  />
        <StatCard title="Blood Managers" value={stats.blood_managers} color="red"    />
        <StatCard title="Patients"       value={stats.patients}       color="purple" />
        <StatCard title="Donors"         value={stats.donors}         color="yellow" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Staff Distribution</h2>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={staffData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" tick={{ fontSize: 12 }} />
              <YAxis allowDecimals={false} />
              <Tooltip />
              <Bar dataKey="value" fill="#3b82f6" radius={[4,4,0,0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Patients vs Donors</h2>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie data={pieData} cx="50%" cy="50%" outerRadius={90} dataKey="value" label>
                {pieData.map((_, i) => <Cell key={i} fill={COLORS[i]} />)}
              </Pie>
              <Legend />
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  )
}
