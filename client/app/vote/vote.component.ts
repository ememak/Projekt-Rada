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
      answer: "",
    },
    {
      question: "Pytanie zamknięte",
      optionsList: [{name: "Opcja 1.", selected: false}, {name: "Opcja 2.", selected: false}],
      type: PollSchema.QuestionType.CLOSE,
      answer: "",
    },
    {
      question: "Pytanie wielokrotnego wyboru",
      optionsList: [{name: "Opcja 1.", selected: false}, {name: "Opcja 2.", selected: false}],
      type: PollSchema.QuestionType.CHECKBOX,
      answer: "",
    },
  ];
  closedOption: number[];

  constructor () {
    this.closedOption = new Array(this.questionsList.length).fill(-1);
  }

  trackOption(index: number, option: string) {
    return index;
  }

  get diagnostic() { return JSON.stringify(this.closedOption/*questionsList*/); }

  onSubmit() {
    for(var i=0; i<this.questionsList.length; i++) {
      var qa = this.questionsList[i];
      if(qa.type == 2) {
        for(var j=0; j<qa.optionsList.length; j++) {
          if(j == this.closedOption[i])
            qa.optionsList[j].selected = true;
          else
            qa.optionsList[j].selected = false;
        }
      }
    }
    console.log(JSON.stringify(this.questionsList));
  }
}


/*
Copyright Google LLC. All Rights Reserved.
Use of this source code is governed by an MIT-style license that
can be found in the LICENSE file at https://angular.io/license
*/
