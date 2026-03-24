import { useState } from 'react'
import { recordDonation } from '../../api/blood'
import { useAuth } from '../../context/AuthContext'
import PageHeader from '../../components/PageHeader'
import Alert from '../../components/Alert'

export default function BloodManagerDonate() {
  const { user } = useAuth()
  const [form,    setForm]    = useState({ donor_email: '', manager_id: user?.user_id || '' })
  const [alert,   setAlert]   = useState({ type: '', message: '' })
  const [loading, setLoading] = useState(false)
  const [result,  setResult]  = useState(null)

  const showAlert = (type, message) => {
    setAlert({ type, message })
    setTimeout(() => setAlert({ type: '', message: '' }), 5000)
  }

  const handleSubmit = async () => {
    if (!form.donor_email) return showAlert('error', 'Donor email is required')
    if (!form.manager_id)  return showAlert('error', 'Manager ID is required')
    setLoading(true)
    setResult(null)
    try {
      const res = await recordDonation({
        donor_email: form.donor_email,
        manager_id:  parseInt(form.manager_id),
      })
      setResult(res.data)
      showAlert('success', 'Blood donation recorded successfully')
      setForm({ ...form, donor_email: '' })
    } catch (err) {
      showAlert('error', err.response?.data?.error || 'Donation failed — transaction rolled back')
    } finally { setLoading(false) }
  }

  return (
    <div>
      <PageHeader title="Record Blood Donation" subtitle="Record a new blood donation — 90 day eligibility is enforced" />

      <div className="max-w-lg">
        <Alert type={alert.type} message={alert.message} onClose={() => setAlert({ type: '', message: '' })} />

        <div className="card space-y-4">
          <div className="p-3 bg-blue-50 rounded-lg text-sm text-blue-700 border border-blue-200">
            ℹ️ Donor must not have donated in the last <strong>90 days</strong>. Blood expiry is auto-set to <strong>42 days</strong> from today.
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Donor Email</label>
            <input
              className="input-field"
              value={form.donor_email}
              onChange={e => setForm({...form, donor_email: e.target.value})}
              placeholder="donor@email.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Manager ID</label>
            <input
              type="number"
              className="input-field"
              value={form.manager_id}
              onChange={e => setForm({...form, manager_id: e.target.value})}
            />
          </div>

          <button onClick={handleSubmit} disabled={loading} className="btn-primary w-full">
            {loading ? 'Recording...' : 'Record Donation'}
          </button>
        </div>

        {result && (
          <div className="card mt-4 bg-green-50 border-green-200">
            <h3 className="font-semibold text-green-800 mb-2"> Donation Recorded</h3>
            <p className="text-sm text-green-700">Donation ID: <strong>#{result.donation_id}</strong></p>
            <p className="text-sm text-green-700">Donor: <strong>{result.donor}</strong></p>
            <p className="text-sm text-green-700">Blood Group: <strong>{result.blood_group}</strong></p>
            <p className="text-sm text-green-700">Blood ID: <strong>#{result.blood_id}</strong></p>
            <p className="text-sm text-green-700">Expiry Date: <strong>{result.expiry_date}</strong></p>
            <p className="text-sm text-green-700">Units Added: <strong>{result.units_added}</strong></p>
          </div>
        )}
      </div>
    </div>
  )
}
