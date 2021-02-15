import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { grpc } from '@improbable-eng/grpc-web';
import { Query } from "Projekt_Rada/query/query_pb_service";
import { PollQuestion, PollSchema } from "Projekt_Rada/query/query_pb";
import { QAListToSchema } from "../proto_parsing";
import { host } from '../host';

@Component({
  selector: 'poll-init',
  templateUrl: './pollinit.component.html',
})
export class PollInitComponent {
  questionsList: PollSchema.QA.AsObject[] = [{
      question: "",
      optionsList: [""],
      type: PollSchema.QuestionType.OPEN,
      answersList: [""],
    },
  ];

  constructor (private router: Router) {}

  addQuestion() {
    this.questionsList.push({
      question: "",
      optionsList: [""],
      type: PollSchema.QuestionType.OPEN,
      answersList: [""],
    });
  }

  addOption(index: number) {
    console.log(index)
    this.questionsList[index].optionsList.push("");
    this.questionsList[index].answersList.push("");
  }

  trackOption(index: number, option: string) {
    return index;
  }

  onSubmit() {
    if (confirm('Czy chcesz wysłać ankietę?')) {
      this.sendPoll();
    }
  }

  sendPoll() {
    const schema: PollSchema = QAListToSchema(this.questionsList);
    grpc.unary(Query.PollInit, {
      request: schema,
      host: host,
      onEnd: res => {
        const { status, statusMessage, headers, message, trailers } = res;
        if (status === grpc.Code.OK && message) {
          let response = (<PollQuestion> message);
          let tokens = response.getTokensList();
          let pollid: number = response.getId();
          this.download("tokeny_" + pollid.toString() + ".txt", tokens);
          this.router.navigate(['/results', pollid]);
        }
      }
    });
  }

  download(filename:string, text: string[]) {
    var element = document.createElement('a');
    element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(text.join("\n")));
    element.setAttribute('download', filename);

    element.style.display = 'none';
    document.body.appendChild(element);

    element.click();
    document.body.removeChild(element);
  }

  get diagnostic() { return JSON.stringify(this.questionsList); }
}

