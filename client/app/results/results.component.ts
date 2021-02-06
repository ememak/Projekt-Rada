import { Component } from '@angular/core';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollSchema, PollSummary } from "Projekt_Rada/query/query_pb";

const host = "http://localhost:12345";

@Component({
  selector: 'results',
  templateUrl: './results.component.html',
})
export class ResultsComponent {
  summary: PollSummary.AsObject = summ;

  input: any[] = [];

  constructor () {
    for (let qa of this.summary.schema.questionsList) {
      let r = [];
      for(let [i, opt] of qa.optionsList.entries()) {
        r.push([opt, parseInt(qa.answersList[i])]);
      }
      this.input.push(r)
    }
    console.log("here")
  }

  trackOption(index: number, option: string) {
    return index;
  }

  get diagnostic() { return JSON.stringify(this.summary.schema); }
  
  onSubmit() {}
}

var summ: PollSummary.AsObject = {
  id: 1,
  votescount: 3,
  schema: {
    questionsList: [{
      question: "Pytanie otwarte",
      optionsList: [],
      type: PollSchema.QuestionType.OPEN,
      answersList: ["Odpowiedź", "Kolejna odpowiedź", "Jeszcze jedna odpowiedź"],
    },
    {
      question: "Pytanie zamknięte",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CLOSE,
      answersList: ["2", "1"],
    },
    {
      question: "Pytanie wielokrotnego wyboru",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CHECKBOX,
      answersList: ["3", "1"],
    },]
  }
}
