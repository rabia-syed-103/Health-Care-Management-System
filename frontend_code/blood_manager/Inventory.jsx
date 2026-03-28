import { useState, useEffect } from 'react'
import { getBloodInventory } from '../../api/blood'
import PageHeader from '../../components/PageHeader'
import Table from '../../components/Table'
import LoadingSpinner from '../../components/LoadingSpinner'
import Alert from '../../components/Alert'

const bloodGroups = ['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-']

export default function BloodManagerInventory() {
  const [inventory, setInventory] = useState([])
  const [filtered,  setFiltered]  = useState([])
  const [loading,   setLoading]   = useState(true)
  const [filterBG,  setFilterBG]  = useState('')
  const [filterUrgency, setFilterUrgency] = useState('')
  const [alert,     setAlert]     = useState({ type: '', message: '' })

  const fetchData = async () => {
    try {
      const res = await getBloodInventory()
      const data = res.data.inventory || []
      setInventory(data)
      setFiltered(data)
    } catch { setInventory([]); setFiltered([]) }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchData() }, [])

  useEffect(() => {
    let result = inventory
    if (filterBG)      result = result.filter(i => i.blood_group?.trim() === filterBG)
    if (filterUrgency) result = result.filter(i => i.expiry_alert === filterUrgency)  // ✅ was i.urgency
    setFiltered(result)
  }, [filterBG, filterUrgency, inventory])

  const urgencyBadge = (alert) => {
    const map = {
      'CRITICAL':      'badge-danger',
      'EXPIRING SOON': 'badge-warning',
      'OK':            'badge-success',
    }
    return <span className={map[alert] || 'badge-info'}>{alert}</span>
  }

  const statusBadge = (status) => {
    const map = { available: 'badge-success', reserved: 'badge-warning' }
    return <span className={map[status] || 'badge-info'}>{status}</span>
  }

  const columns = [
    { key: 'id',           label: 'Blood ID' },
    { key: 'blood_group',  label: 'Blood Group',  render: (row) => <span className="badge-danger">{row.blood_group?.trim()}</span> },
    { key: 'units',        label: 'Units' },
    { key: 'status',       label: 'Status',       render: (row) => statusBadge(row.status) },
    { key: 'expiry_date',  label: 'Expiry Date' },
    { key: 'expiry_alert', label: 'Urgency',      render: (row) => urgencyBadge(row.expiry_alert) },  // ✅ was row.urgency
  ]

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Blood Inventory"
        subtitle={`${filtered.length} of ${inventory.length} blood units`}
      />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

      <div className="flex gap-3 mb-4 flex-wrap">
        <select className="input-field max-w-xs" value={filterBG} onChange={e => setFilterBG(e.target.value)}>
          <option value="">All Blood Groups</option>
          {bloodGroups.map(bg => <option key={bg} value={bg}>{bg}</option>)}
        </select>
        <select className="input-field max-w-xs" value={filterUrgency} onChange={e => setFilterUrgency(e.target.value)}>
          <option value="">All Urgency Levels</option>
          <option value="CRITICAL">Critical</option>
          <option value="EXPIRING SOON">Expiring Soon</option>
          <option value="OK">OK</option>
        </select>
        {(filterBG || filterUrgency) && (
          <button onClick={() => { setFilterBG(''); setFilterUrgency('') }} className="btn-secondary">Clear</button>
        )}
      </div>

      <Table columns={columns} data={filtered} emptyMessage="No blood inventory found" />
    </div>
  )
}