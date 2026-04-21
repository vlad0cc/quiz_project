import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { createSession, getResult } from '../api/client';
import { ProfileRadar } from '../components/ProfileRadar';
import type { Result } from '../types';
import styles from './ResultPage.module.css';

const SESSION_STORAGE_KEY = 'profil-math-session-id';

const PROFILE_DESCRIPTIONS: Record<'AD' | 'ZI', string> = {
  AD: 'Профиль подходит тем, кому интересны закономерности, статистика, модели и работа с неопределённостью в данных.',
  ZI: 'Профиль подходит тем, кому ближе логические правила, контроль доступа, дискретные структуры и надёжность систем.',
};

export function ResultPage() {
  const navigate = useNavigate();
  const { sessionID } = useParams<{ sessionID: string }>();
  const [result, setResult] = useState<Result | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!sessionID) {
      navigate('/', { replace: true });
      return;
    }

    void loadResult(sessionID);
  }, [navigate, sessionID]);

  async function loadResult(currentSessionID: string) {
    setLoading(true);
    setError(null);

    try {
      const data = await getResult(currentSessionID);
      setResult(data);
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось загрузить результат');
    } finally {
      setLoading(false);
    }
  }

  async function handleRestart() {
    setError(null);
    try {
      const response = await createSession();
      sessionStorage.setItem(SESSION_STORAGE_KEY, response.session_id);
      navigate('/quiz');
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось создать новую сессию');
    }
  }

  if (loading || !result) {
    return (
      <main className={styles.page}>
        <section className={styles.card}>
          <p className={styles.state}>Формируем итоговую рекомендацию...</p>
          {error ? <p className={styles.error}>{error}</p> : null}
        </section>
      </main>
    );
  }

  return (
    <main className={styles.page}>
      <section className={styles.card}>
        <span className={styles.badge}>Результат диагностики</span>
        <h1 className={styles.title}>{result.recommendation_label}</h1>
        <p className={styles.subtitle}>
          Уверенность рекомендации: <strong>{Math.round(result.confidence * 100)}%</strong>
        </p>

        <ProfileRadar
          adPercent={result.breakdown.ad_percent}
          ziPercent={result.breakdown.zi_percent}
        />

        <div className={styles.scores}>
          <div>
            <span>Баллы AD</span>
            <strong>{result.score_ad.toFixed(1)}</strong>
          </div>
          <div>
            <span>Баллы ZI</span>
            <strong>{result.score_zi.toFixed(1)}</strong>
          </div>
        </div>

        <div className={styles.profiles}>
          <article>
            <h2>Анализ данных</h2>
            <p>{PROFILE_DESCRIPTIONS.AD}</p>
          </article>
          <article>
            <h2>Защита информации</h2>
            <p>{PROFILE_DESCRIPTIONS.ZI}</p>
          </article>
        </div>

        <button className={styles.button} type="button" onClick={handleRestart}>
          Пройти ещё раз
        </button>
        {error ? <p className={styles.error}>{error}</p> : null}
      </section>
    </main>
  );
}
