export interface Question {
  question_id: number;
  text: string;
  options: Record<string, string>;
  step: number;
  total_steps: number;
}

export interface AnswerResult {
  is_correct: boolean;
  step: number;
  total_steps: number;
}

export interface Result {
  recommendation: 'AD' | 'ZI';
  recommendation_label: string;
  confidence: number;
  score_ad: number;
  score_zi: number;
  breakdown: {
    ad_percent: number;
    zi_percent: number;
  };
}
