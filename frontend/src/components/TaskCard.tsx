import styles from './TaskCard.module.css';

interface TaskCardProps {
  title: string;
  text: string;
  options: Record<string, string>;
  selectedOption: string | null;
  disabled: boolean;
  onSelect: (option: string) => void;
}

export function TaskCard({
  title,
  text,
  options,
  selectedOption,
  disabled,
  onSelect,
}: TaskCardProps) {
  return (
    <section className={styles.card}>
      <div className={styles.header}>
        <span className={styles.badge}>{title}</span>
        <h2 className={styles.title}>{text}</h2>
      </div>

      <div className={styles.options}>
        {Object.entries(options).map(([key, value]) => {
          const isSelected = selectedOption === key;
          return (
            <button
              key={key}
              type="button"
              className={`${styles.option} ${isSelected ? styles.optionSelected : ''}`}
              disabled={disabled}
              onClick={() => onSelect(key)}
            >
              <span className={styles.optionKey}>{key}</span>
              <span>{value}</span>
            </button>
          );
        })}
      </div>
    </section>
  );
}
