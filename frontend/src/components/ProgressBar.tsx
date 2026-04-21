import styles from './ProgressBar.module.css';

interface ProgressBarProps {
  step: number;
  totalSteps: number;
}

export function ProgressBar({ step, totalSteps }: ProgressBarProps) {
  const progress = Math.min((step / totalSteps) * 100, 100);

  return (
    <div className={styles.wrapper}>
      <div className={styles.meta}>
        <span>Диагностика</span>
        <span>
          Шаг {step} из {totalSteps}
        </span>
      </div>
      <div className={styles.track}>
        <div className={styles.value} style={{ width: `${progress}%` }} />
      </div>
    </div>
  );
}
