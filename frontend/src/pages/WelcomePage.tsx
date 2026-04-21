import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { createSession } from '../api/client';
import styles from './WelcomePage.module.css';

const SESSION_STORAGE_KEY = 'profil-math-session-id';

export function WelcomePage() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleStart() {
    setLoading(true);
    setError(null);

    try {
      const response = await createSession();
      sessionStorage.setItem(SESSION_STORAGE_KEY, response.session_id);
      navigate('/quiz');
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось создать сессию');
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className={styles.page}>
      <section className={styles.hero}>
        <span className={styles.kicker}>Профиль.Математика</span>
        <h1 className={styles.title}>Диагностика, которая переводит абстрактный выбор профиля в измеримый результат.</h1>
        <p className={styles.description}>
          Платформа сравнивает ваши предпочтения в задачах по логике, структурам данных, вероятности и анализу
          закономерностей. На выходе вы получаете рекомендацию между двумя профилями кафедры прикладной математики.
        </p>

        <div className={styles.metrics}>
          <div>
            <strong>10 задач</strong>
            <span>адаптивный маршрут без регистрации</span>
          </div>
          <div>
            <strong>~7 минут</strong>
            <span>полный цикл до результата</span>
          </div>
          <div>
            <strong>2 профиля</strong>
            <span>Анализ данных и Защита информации</span>
          </div>
        </div>

        <button className={styles.button} type="button" disabled={loading} onClick={handleStart}>
          {loading ? 'Создаём сессию...' : 'Начать диагностику'}
        </button>

        {error ? <p className={styles.error}>{error}</p> : null}
      </section>
    </main>
  );
}
