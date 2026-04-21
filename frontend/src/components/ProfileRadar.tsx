import styles from './ProfileRadar.module.css';

interface ProfileRadarProps {
  adPercent: number;
  ziPercent: number;
}

export function ProfileRadar({ adPercent, ziPercent }: ProfileRadarProps) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.row}>
        <div className={styles.labelBlock}>
          <span className={`${styles.dot} ${styles.adDot}`} />
          <span>Анализ данных</span>
        </div>
        <strong>{adPercent}%</strong>
      </div>
      <div className={styles.bar}>
        <div className={`${styles.fill} ${styles.adFill}`} style={{ width: `${adPercent}%` }} />
      </div>

      <div className={styles.row}>
        <div className={styles.labelBlock}>
          <span className={`${styles.dot} ${styles.ziDot}`} />
          <span>Защита информации</span>
        </div>
        <strong>{ziPercent}%</strong>
      </div>
      <div className={styles.bar}>
        <div className={`${styles.fill} ${styles.ziFill}`} style={{ width: `${ziPercent}%` }} />
      </div>
    </div>
  );
}
