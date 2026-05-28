import { useEffect, useState } from 'react'

type ThemeMode = 'light'

function applyThemeMode(mode: ThemeMode) {
	document.documentElement.classList.remove('light', 'dark')
	document.documentElement.classList.add('light')
	document.documentElement.setAttribute('data-theme', 'light')
	document.documentElement.style.colorScheme = 'light'
}

export default function ThemeToggle() {
  const [mode, setMode] = useState<ThemeMode>('auto')

  useEffect(() => {
		setMode('light')
		applyThemeMode('light')
	}, [])

	useEffect(() => {
		applyThemeMode('light')
	}, [mode])

	function toggleMode() {
		setMode('light')
		applyThemeMode('light')
		window.localStorage.setItem('theme', 'light')
	}

	const label = 'Theme mode: light.'

  return (
    <button
      type="button"
      onClick={toggleMode}
      aria-label={label}
      title={label}
      className="rounded-full border border-[var(--chip-line)] bg-[var(--chip-bg)] px-3 py-1.5 text-sm font-semibold text-[var(--sea-ink)] shadow-[0_8px_22px_rgba(30,90,72,0.08)] transition hover:-translate-y-0.5"
    >
		Light
	</button>
  )
}
