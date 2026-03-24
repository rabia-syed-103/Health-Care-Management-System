import { useState, useEffect } from 'react'
import { getMedicineStock } from '../../api/medicines'
import { getPendingPrescriptions } from '../../api/prescriptions'
import PageHeader from '../../components/PageHeader'
import StatCard from '../../components/StatCard'
import LoadingSpinner from '../../components/LoadingSpinner'
import { useAuth } from '../../context/AuthContext'
import { PieChart, Pie, Cell, Legend, Tooltip, ResponsiveContainer } from 'recharts'

export default function PharmacistDashboard() {
  const { user } = useAuth()
  const [stats,   setStats]   = useState({ total: 0, outOfStock: 0, expiring: 0, pending: 0 })
  const [pieData, setPieData] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchAll = async () => {
      try {
        const [stock, pending] = await Promise.all([
          getMedicineStock(),
          getPendingPrescriptions(),
        ])
        const medicines = stock.data.medicines || []
        const outOfStock = medicines.filter(m => m.status === 'OUT OF STOCK').length
        const expired    = medicines.filter(m => m.status === 'EXPIRED').length
        const critical   = medicines.filter(m => m.status === 'CRITICAL').length
        const expiring   = medicines.filter(m => m.status === 'EXPIRING SOON').length
        const ok         = medicines.filter(m => m.status === 'OK').length

        setStats({
          total:      medicines.length,
          outOfStock: outOfStock + expired,
          expiring:   critical + expiring,
          pending:    pending.data.prescriptions?.length || 0,
        })

        setPieData([
          { name: 'OK',           value: ok },
          { name: 'Expiring Soon',value: expiring },
          { name: 'Critical',     value: critical },
          { name: 'Out of Stock', value: outOfStock },
          { name: 'Expired',      value: expired },
        ].filter(d => d.value > 0))

      } catch { }
      finally { setLoading(false) }
    }
    fetchAll()
  }, [])

  const COLORS = ['#22c55e', '#f59e0b', '#ef4444', '#6b7280', '#dc2626']

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title={`Welcome, ${user?.name || 'Pharmacist'}`}
        subtitle="Medicine inventory and prescription overview"
      />

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <StatCard title="Total Medicines"    value={stats.total}      color="blue"   />
        <StatCard title="Out of Stock"       value={stats.outOfStock} color="red"    />
        <StatCard title="Expiring / Critical"value={stats.expiring}   color="yellow" />
        <StatCard title="Pending Rx"         value={stats.pending}    color="purple" />
      </div>

      {pieData.length > 0 && (
        <div className="card max-w-lg">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Medicine Stock Status</h2>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie data={pieData} cx="50%" cy="50%" outerRadius={90} dataKey="value" label>
                {pieData.map((_, i) => <Cell key={i} fill={COLORS[i % COLORS.length]} />)}
              </Pie>
              <Legend />
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  )
}
