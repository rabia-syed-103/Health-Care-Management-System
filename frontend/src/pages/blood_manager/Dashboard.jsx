import { useState, useEffect } from 'react'
import { getDonationHistory, getBloodInventory, getPendingRequests, getExpiredBlood } from '../../api/blood'
import { getAllDonors } from '../../api/donors'
import PageHeader from '../../components/PageHeader'
import StatCard from '../../components/StatCard'
import LoadingSpinner from '../../components/LoadingSpinner'
import { useAuth } from '../../context/AuthContext'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from 'recharts'

export default function BloodManagerDashboard() {
  const { user } = useAuth()
  const [stats,     setStats]     = useState({ donors: 0, donations: 0, pending: 0, expired: 0 })
  const [inventory, setInventory] = useState([])
  const [loading,   setLoading]   = useState(true)

  useEffect(() => {
    const fetchAll = async () => {
      try {
        const [donors, donations, pending, expired, inv] = await Promise.all([
          getAllDonors(),
          getDonationHistory(),
          getPendingRequests(),
          getExpiredBlood(),
          getBloodInventory(),
        ])
        setStats({
        donors:    donors.data.donors?.length             || 0,
        donations: donations.data.donations?.length       || 0,
        pending:   pending.data.pending_requests?.length  || 0,
        expired:   expired.data.expired_blood?.length     || 0,
      })
        setInventory(inv.data.inventory || [])
      } catch { }
      finally { setLoading(false) }
    }
    fetchAll()
  }, [])

  // Group inventory by blood group for chart
  const chartData = inventory.reduce((acc, item) => {
    const existing = acc.find(a => a.group === item.blood_group)
    if (existing) existing.units += item.units
    else acc.push({ group: item.blood_group, units: item.units })
    return acc
  }, []).sort((a, b) => a.group.localeCompare(b.group))

  const COLORS = ['#ef4444','#f97316','#eab308','#22c55e','#14b8a6','#3b82f6','#8b5cf6','#ec4899']

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title={`Welcome, ${user?.email  || 'Blood Manager'}`}
        badge={user?.id ? `ID: ${user.id}` : null}
        subtitle="Blood bank overview and inventory status"
      />

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <StatCard title="Total Donors"    value={stats.donors}    color="blue"   />
        <StatCard title="Total Donations" value={stats.donations} color="green"  />
        <StatCard title="Pending Requests"value={stats.pending}   color="yellow" />
        <StatCard title="Expired Units"   value={stats.expired}   color="red"    />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Blood Inventory by Group</h2>
          {chartData.length === 0 ? (
            <p className="text-gray-400 text-sm text-center py-8">No inventory data</p>
          ) : (
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="group" tick={{ fontSize: 12 }} />
                <YAxis allowDecimals={false} />
                <Tooltip />
                <Bar dataKey="units" fill="#ef4444" radius={[4,4,0,0]} />
              </BarChart>
            </ResponsiveContainer>
          )}
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Blood Group Distribution</h2>
          {chartData.length === 0 ? (
            <p className="text-gray-400 text-sm text-center py-8">No inventory data</p>
          ) : (
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie data={chartData} dataKey="units" nameKey="group" cx="50%" cy="50%" outerRadius={90} label>
                  {chartData.map((_, i) => <Cell key={i} fill={COLORS[i % COLORS.length]} />)}
                </Pie>
                <Legend />
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          )}
        </div>
      </div>
    </div>
  )
}
