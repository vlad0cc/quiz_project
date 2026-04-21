import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { getNextQuestion, submitAnswer } from '../api/client';
import { ProgressBar } from '../components/ProgressBar';
import { TaskCard } from '../components/TaskCard';
import type { Question } from '../types';
import styles from './QuestionPage.module.css';

const SESSION_STORAGE_KEY = 'profil-math-session-id';
const ANSWER_DELAY_MS = 600;

function wait(delay: number) {
  return new Promise((resolve) => {
    window.setTimeout(resolve, delay);
  });
}

export function QuestionPage() {
  const navigate = useNavigate();
  const [question, setQuestion] = useState<Question | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [selectedOption, setSelectedOption] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const sessionId = sessionStorage.getItem(SESSION_STORAGE_KEY);

  useEffect(() => {
    if (!sessionId) {
      navigate('/', { replace: true });
      return;
    }

    void loadQuestion(sessionId);
  }, [navigate, sessionId]);

  async function loadQuestion(currentSessionId: string) {
    setLoading(true);
    setSelectedOption(null);
    setError(null);

    try {
      const nextQuestion = await getNextQuestion(currentSessionId);
      if (!nextQuestion) {
        navigate(`/result/${currentSessionId}`, { replace: true });
        return;
      }
      setQuestion(nextQuestion);
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось загрузить вопрос');
    } finally {
      setLoading(false);
    }
  }

  async function handleSelect(option: string) {
    if (!sessionId || !question || submitting) {
      return;
    }

    setSubmitting(true);
    setSelectedOption(option);
    setError(null);

    try {
      await submitAnswer(sessionId, question.question_id, option);
      await wait(ANSWER_DELAY_MS);
      await loadQuestion(sessionId);
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось отправить ответ');
    } finally {
      setSubmitting(false);
    }
  }

  if (loading || !question) {
    return (
      <main className={styles.page}>
        <section className={styles.panel}>
          <p className={styles.state}>Загружаем следующую задачу...</p>
          {error ? <p className={styles.error}>{error}</p> : null}
        </section>
      </main>
    );
  }

  return (
    <main className={styles.page}>
      <section className={styles.panel}>
        <ProgressBar step={question.step} totalSteps={question.total_steps} />
        <TaskCard
          title="Ситуационная задача"
          text={question.text}
          options={question.options}
          selectedOption={selectedOption}
          disabled={submitting}
          onSelect={handleSelect}
        />
        {error ? <p className={styles.error}>{error}</p> : null}
      </section>
    </main>
  );
}
