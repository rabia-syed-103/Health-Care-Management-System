import { useState, useEffect } from 'react'
import { getExpiredBlood } from '../../api/blood'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import LoadingSpinner from '../../components/LoadingSpinner'
import Alert from '../../components/Alert'

export default function BloodManagerExpired() {
  const [expired,  setExpired]  = useState([])
  const [filtered, setFiltered] = useState([])
  const [loading,  setLoading]  = useState(true)
  const [filterBG, setFilterBG] = useState('')
  const [alert,    setAlert]    = useState({ type: '', message: '' })

  const bloodGroups = ['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-']

const fetchData = async () => {
  try {
    const res = await getExpiredBlood()
    const data = res.data.expired_blood?.map(b => ({
      ...b,
      blood_group: b.blood_group?.trim()
    })) || []
    setExpired(data)
    setFiltered(data)
  } catch {
    setExpired([])
    setFiltered([])
  } finally {
    setLoading(false)
  }
}

  useEffect(() => { fetchData() }, [])

  useEffect(() => {
    let result = expired
    if (filterBG) result = result.filter(e => e.blood_group?.trim() === filterBG)
    setFiltered(result)
  }, [filterBG, expired])

const columns = [
  { key: 'id',         label: 'Blood ID' },
  { key: 'blood_group', label: 'Blood Group', render: (row) => <span className="badge-danger">{row.blood_group?.trim()}</span> },
  { key: 'units',       label: 'Units' },
  { key: 'expiry_date', label: 'Expiry Date' },
  { key: 'days_expired',label: 'Days Expired', render: (row) => (
    <span className="text-red-600 font-medium">{row.days_expired} days</span>
  )},
]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Expired Blood"
        subtitle={`${filtered.length} expired blood units`}
      />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

      <div className="flex gap-3 mb-4">
        <select className="input-field max-w-xs" value={filterBG} onChange={e => setFilterBG(e.target.value)}>
          <option value="">All Blood Groups</option>
          {bloodGroups.map(bg => <option key={bg} value={bg}>{bg}</option>)}
        </select>
        {filterBG && (
          <button onClick={() => setFilterBG('')} className="btn-secondary">Clear</button>
        )}
      </div>

      {filtered.length === 0 ? (
        <div className="card text-center py-12">
          <p className="text-4xl mb-3">✅</p>
          <p className="text-gray-500">No expired blood units found</p>
        </div>
      ) : (
        <Table columns={columns} data={filtered} emptyMessage="No expired blood found" />
      )}
    </div>
  )
}
