import { OnInit, Component, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs';
import { MessageService } from 'primeng/api';
import { SearchBinaryService, CheckNames, FileData } from './search-binary.service';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {

  @ViewChild('scanForm') scanForm: NgForm | undefined;

  title = 'search-binary';
  display: boolean = false;
  scanInProgress: boolean = false;
  path: string = '';
  placeholder: string = '';
  selectedCheckNames: CheckNames[] = [];
  selectedCheckNamesShort: string[] = [];
  checknames: CheckNames[] = [
    {name: 'Binary files', key: 'SearchBinaries'},
    {name: 'Committer name or email does not match criteria', key: 'CheckCommitMetaInformation'},
    {name: 'Illegal unicode characters', key: 'SearchIllegalUnicodeCharacters'},
    {name: 'Large files', key: 'SearchBigFiles'}
  ];
  getDataSubscription: Subscription | undefined;
  fileData: FileData[] = [];

  constructor(
    private messageService: MessageService,
    private searchBinaryService: SearchBinaryService) {
  }

  ngOnInit() {
  }

  /**
   * Open the scan config dialog.
   */
  openConfigDialog() {
    //this.fileData = [];
    //this.selectedCheckNames = [];
    this.display = true;
    this.path = '';
    this.placeholder = 'https://github.com/grafana/loki,git@github.com/UserX/gitscanner.git,file://C:\\Users\\UserX\\myRepo';
    this.scanForm?.form.markAsPristine();
  }

  /**
   * Empty the input field placeholder once in focus.
   */
  emptyPlaceholder() {
    this.placeholder = '';
  }

  /**
   * Submit the scan config and perform scan.
   */
  submit() {
    this.fileData = [];
    let selected: string[] = [];
    this.selectedCheckNames.forEach((entry) => {
      selected.push(entry.key);
    });
    this.selectedCheckNamesShort = Object.assign([], selected);
    this.getDataSubscription = this.searchBinaryService.getFileData(this.path, selected).subscribe(data => {
      if (data && data.length > 0) {
        this.fileData = data;
      }
      this.scanInProgress = false;
     },
     error => {
      this.messageService.add({key: 'scanError', severity:'error', summary:'Error', detail: error.statusText});
      this.scanInProgress = false;
     });
    this.display = false;
    this.scanInProgress = true;
  }

  ngOnDestroy() {
    this.getDataSubscription?.unsubscribe();
  }
}
