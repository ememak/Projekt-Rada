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
  answers: PollSummary.AsObject = poll;

  constructor () {
    for(let question of this.answers.schema.questionsList){
      if(question.type != PollSchema.QuestionType.OPEN){
        for(let [j, ans] of question.answersList.entries()){
          if(ans == ""){
            question.answersList[j] = "0"
          }
        }
      }
      else {
        question.answersList = []
      }
    }
    for(let vote of this.answers.votesList){
      for(let [i, question] of vote.questionsList.entries()){
        if(question.type == PollSchema.QuestionType.OPEN){
          this.answers.schema.questionsList[i].answersList.push(question.answersList[0])
        }
        if(vote.questionsList[i].type == PollSchema.QuestionType.CLOSE){
          for(let [j, ans] of question.answersList.entries()){
            let v = parseInt(this.answers.schema.questionsList[i].answersList[j]);
            this.answers.schema.questionsList[i].answersList[j] = (v + parseInt(ans)).toString()
          }
        }
        if(vote.questionsList[i].type == PollSchema.QuestionType.CHECKBOX){
          for(let [j, ans] of question.answersList.entries()){
            let v = parseInt(this.answers.schema.questionsList[i].answersList[j]);
            this.answers.schema.questionsList[i].answersList[j] = (v + parseInt(ans)).toString()
          }
        }
      }
    }
  }

  addOption(index: number) {
    console.log(index)
  }

  trackOption(index: number, option: string) {
    return index;
  }

  get diagnostic() { return JSON.stringify(this.answers.schema); }
  
  onSubmit() {}
}

var poll: PollSummary.AsObject = {
  id: 1,
  schema: {
    questionsList: [{
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
    },]
  },
  votesList: [{
    questionsList: [{
      question: "Pytanie otwarte",
      optionsList: [],
      type: PollSchema.QuestionType.OPEN,
      answersList: ["Odpowiedź"],
    },
    {
      question: "Pytanie zamknięte",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CLOSE,
      answersList: ["1", "0"],
    },
    {
      question: "Pytanie wielokrotnego wyboru",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CHECKBOX,
      answersList: ["1", "1"],
    },]
  },
  {
    questionsList: [{
      question: "Pytanie otwarte",
      optionsList: [],
      type: PollSchema.QuestionType.OPEN,
      answersList: ["Kolejna odpowiedź"],
    },
    {
      question: "Pytanie zamknięte",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CLOSE,
      answersList: ["0", "1"],
    },
    {
      question: "Pytanie wielokrotnego wyboru",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CHECKBOX,
      answersList: ["1", "0"],
    }]
  },
  {
    questionsList: [{
      question: "Pytanie otwarte",
      optionsList: [],
      type: PollSchema.QuestionType.OPEN,
      answersList: ["Jeszcze jedna odpowiedź"],
    },
    {
      question: "Pytanie zamknięte",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CLOSE,
      answersList: ["1", "0"],
    },
    {
      question: "Pytanie wielokrotnego wyboru",
      optionsList: ["Opcja 1.", "Opcja 2."],
      type: PollSchema.QuestionType.CHECKBOX,
      answersList: ["1", "0"],
    }]
  },
  ]
}
