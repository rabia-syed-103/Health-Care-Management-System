import { useState, useEffect } from 'react'
import { getPendingPrescriptions } from '../../api/prescriptions'
import PageHeader from '../../components/PageHeader'
import LoadingSpinner from '../../components/LoadingSpinner'
import Alert from '../../components/Alert'

export default function PharmacistPending() {
  const [prescriptions, setPrescriptions] = useState([])
  const [loading,       setLoading]       = useState(true)
  const [alert,         setAlert]         = useState({ type: '', message: '' })

  useEffect(() => {
    const fetch = async () => {
      try {
        const res = await getPendingPrescriptions()
        setPrescriptions(res.data.prescriptions || [])
      } catch {
        setAlert({ type: 'error', message: 'Failed to load prescriptions' })
      } finally { setLoading(false) }
    }
    fetch()
  }, [])

  if (loading) return <LoadingSpinner />

  return (
    <div>
      <PageHeader
        title="Pending Prescriptions"
        subtitle={`${prescriptions.length} prescriptions awaiting dispensing`}
      />
      <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

      {prescriptions.length === 0 ? (
        <div className="card text-center py-12 text-gray-400">No pending prescriptions</div>
      ) : (
        <div className="space-y-4">
          {prescriptions.map((rx, i) => (
            <div key={i} className="card">
              <div className="flex items-start justify-between mb-3">
                <div>
                  <p className="font-semibold text-gray-800">Prescription #{rx.prescription_id}</p>
                  <p className="text-sm text-gray-500">{rx.date}</p>
                </div>
                <span className="badge-warning">Pending</span>
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm mb-3">
                <div><span className="text-gray-500">Patient:</span> <strong>{rx.patient_name}</strong></div>
                <div><span className="text-gray-500">Doctor:</span> <strong>{rx.doctor_name}</strong></div>
              </div>
              <div>
                <p className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">Medicines</p>
                <div className="space-y-1">
                  {rx.medicines?.map((m, j) => (
                    <div key={j} className="flex justify-between text-sm bg-gray-50 px-3 py-1.5 rounded">
                      <span>{m.name}</span>
                      <span className="text-gray-500">× {m.quantity}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
