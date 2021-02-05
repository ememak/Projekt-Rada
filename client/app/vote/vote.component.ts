import { Component } from '@angular/core';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollSchema } from "Projekt_Rada/query/query_pb";

const host = "http://localhost:12345";

@Component({
  selector: 'vote',
  templateUrl: './vote.component.html',
})
export class VoteComponent {
  questionsList: PollSchema.QA.AsObject[]  = [{
      question: "Pytanie otwarte",
      optionsList: [],
      type: PollSchema.QuestionType.OPEN,
      answersList: [""],
    },
    {
      question: "Pytanie zamknięte",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CLOSE,
      answersList: ["", ""],
    },
    {
      question: "Pytanie wielokrotnego wyboru",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CHECKBOX,
      answersList: ["", ""],
    },
  ];

  constructor () {
  }

  trackOption(index: number, option: string) {
    return index;
  }

  get diagnostic() { return JSON.stringify(this.questionsList); }

  onSubmit() {
    for(let qa of this.questionsList){
      if(qa.type==PollSchema.QuestionType.CLOSE){
        qa.answersList = qa.answersList.map((ans, index) => {
          return parseInt(qa.answersList[0])==index?"true":"false";
        })
      }
      if(qa.type==PollSchema.QuestionType.CHECKBOX){
        qa.answersList = qa.answersList.map((ans, index) => {
          return ans==""?"false":ans;
        })
      }
    }
    console.log(JSON.stringify(this.questionsList));
  }
}

