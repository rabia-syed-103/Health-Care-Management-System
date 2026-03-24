import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const navItems = {
  admin: [
    { label: 'Dashboard',      path: '/admin' },
    { label: 'Doctors',        path: '/admin/doctors' },
    { label: 'Receptionists',  path: '/admin/receptionists' },
    { label: 'Pharmacists',    path: '/admin/pharmacists' },
    { label: 'Blood Managers', path: '/admin/blood-managers' },
    { label: 'Medicine Stock', path: '/admin/medicines' },
  ],
  doctor: [
    { label: 'My Appointments',  path: '/doctor' },
    { label: 'Prescribe',        path: '/doctor/prescribe' },
    { label: 'Blood Request',    path: '/doctor/blood-request' },
    { label: 'Patient History',  path: '/doctor/patient-history' },
  ],
  receptionist: [
    { label: 'Dashboard',      path: '/receptionist' },
    { label: 'Patients',       path: '/receptionist/patients' },
    { label: 'Appointments',   path: '/receptionist/appointments' },
    { label: 'ER Patients',    path: '/receptionist/er-patients' },
    { label: 'ER Entries',     path: '/receptionist/er-entries' },
  ],
  pharmacist: [
    { label: 'Dashboard',          path: '/pharmacist' },
    { label: 'Pending Rx',         path: '/pharmacist/pending' },
    { label: 'Medicine Stock',     path: '/pharmacist/stock' },
    { label: 'Add Medicine',       path: '/pharmacist/add-medicine' },
    { label: 'Dispense',           path: '/pharmacist/dispense' },
  ],
  blood_manager: [
    { label: 'Dashboard',         path: '/blood-manager' },
    { label: 'Donors',            path: '/blood-manager/donors' },
    { label: 'Record Donation',   path: '/blood-manager/donate' },
    { label: 'Blood Inventory',   path: '/blood-manager/inventory' },
    { label: 'Pending Requests',  path: '/blood-manager/requests' },
    { label: 'Expired Blood',     path: '/blood-manager/expired' },
  ],
}

const roleColors = {
  admin:         'from-purple-700 to-purple-900',
  doctor:        'from-blue-700 to-blue-900',
  receptionist:  'from-teal-700 to-teal-900',
  pharmacist:    'from-green-700 to-green-900',
  blood_manager: 'from-red-700 to-red-900',
}

const roleLabels = {
  admin:         'Admin Panel',
  doctor:        'Doctor Portal',
  receptionist:  'Receptionist',
  pharmacist:    'Pharmacist',
  blood_manager: 'Blood Manager',
}

export default function Sidebar() {
  const { role, user, logout } = useAuth()
  const navigate = useNavigate()
  const items = navItems[role] || []
  const gradient = roleColors[role] || 'from-gray-700 to-gray-900'

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className={`w-64 min-h-screen bg-gradient-to-b ${gradient} flex flex-col`}>
      {/* Header */}
      <div className="p-6 border-b border-white/10">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-white/20 flex items-center justify-center text-white font-bold text-lg">
            {user?.name?.charAt(0) || role?.charAt(0).toUpperCase()}
          </div>
          <div>
            <p className="text-white font-semibold text-sm">{user?.name || 'User'}</p>
            <p className="text-white/60 text-xs capitalize">{roleLabels[role]}</p>
          </div>
        </div>
      </div>

      {/* Nav Items */}
      <nav className="flex-1 p-4 space-y-1">
        {items.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            end={item.path === '/admin' || item.path === '/doctor' || item.path === '/receptionist' || item.path === '/pharmacist' || item.path === '/blood-manager'}
            className={({ isActive }) =>
              `block px-4 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 ${
                isActive
                  ? 'bg-white/20 text-white'
                  : 'text-white/70 hover:bg-white/10 hover:text-white'
              }`
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>

      {/* Logout */}
      <div className="p-4 border-t border-white/10">
        <button
          onClick={handleLogout}
          className="w-full px-4 py-2.5 rounded-lg text-sm font-medium text-white/70 hover:bg-white/10 hover:text-white transition-all duration-200 text-left"
        >
          Logout
        </button>
      </div>
    </div>
  )
}