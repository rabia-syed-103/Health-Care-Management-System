export default function PageHeader({ title, subtitle, action, badge }) {
  return (
    <div className="flex items-center justify-between mb-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">{title}</h1>
        {subtitle && <p className="text-gray-500 mt-1 text-sm">{subtitle}</p>}
        {badge && (
          <span className="inline-block mt-2 px-3 py-1 text-xs font-semibold bg-blue-100 text-blue-700 rounded-full">
            {badge}
          </span>
        )}
      </div>
      {action && <div>{action}</div>}
    </div>
  )
}
