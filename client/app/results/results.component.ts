import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollSchema, PollSummary, SummaryRequest } from "Projekt_Rada/query/query_pb";
import { host } from '../host';

@Component({
  selector: 'results',
  templateUrl: './results.component.html',
})
export class ResultsComponent {
  summary: PollSummary.AsObject = summ;

  input: any[] = [];

  pollid: number;

  constructor (private route: ActivatedRoute) {
    this.pollid = parseInt(this.route.snapshot.paramMap.get('pollid'));
    console.log(this.pollid)
    if(isNaN(this.pollid)){
      for (let qa of this.summary.schema.questionsList) {
        let r = [];
        for(let [i, opt] of qa.optionsList.entries()) {
          r.push([opt, parseInt(qa.answersList[i])]);
        }
        this.input.push(r)
      }
    } 
    else {
    let request: SummaryRequest = new SummaryRequest();
    request.setPollid(this.pollid);
    grpc.unary(Query.GetSummary, {
      request: request,
      host: host,
      onEnd: res => {
        const { status, statusMessage, headers, message, trailers } = res;
        console.log("pollInit.onEnd.status", status, statusMessage);
        console.log("pollInit.onEnd.headers", headers);
        if (status === grpc.Code.OK && message) {
          console.log("pollInit.onEnd.message", message.toObject());
          this.summary = (<PollSummary>message).toObject()
          for (let qa of this.summary.schema.questionsList) {
            let r = [];
            for(let [i, opt] of qa.optionsList.entries()) {
              r.push([opt, parseInt(qa.answersList[i])]);
            }
            this.input.push(r)
          }
        }
        console.log("pollInit.onEnd.trailers", trailers);
      }
    });}
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
