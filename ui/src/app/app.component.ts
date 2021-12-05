import { OnInit, Component, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { MessageService } from 'primeng/api';
import { SearchBinaryService, CheckNames, FileData } from './search-binary.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {
  title = 'search-binary';
  display: boolean = false;
  path: string = '';
  selectedCheckNames: string[] = [];
  selectedCheckNamesShort: string[] = [];
  checknames: CheckNames[] = [
    {name: 'Binary files', key: 'github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries'},
    {name: 'Illegal unicode characters', key: 'github.com/MrWong99/gitscanner/checks/unicode.SearchUnicode'},
    {name: 'Commiter email does not match configured pattern', key: 'github.com/MrWong99/gitscanner/checks/commitmeta.CheckCommitAuthor'}
  ];
  getDataSubscription: Subscription | undefined;

  fileData: FileData[] = [
    {
      "date": "2021-12-02T16:16:07.982032+01:00",
      "repository": "git@github.com:MrWong99/micasuca.git",
      "error": "",
      "checks": [
        {
          "origin": "app/src/main/ic_message-playstore.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "12.5 kB" }
        },
        {
          "origin": "app/src/main/ic_star-playstore.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "28.0 kB" }
        },
        {
          "origin": "app/src/main/ic_wishlist-playstore.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "20.9 kB" }
        },
        {
          "origin": "app/src/main/res/drawable-v24/micasuca_view_background.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "27.7 kB" }
        },
        {
          "origin": "app/src/main/res/drawable-v24/micasuca_view_home_cut.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "30.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_launcher.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_launcher_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_message.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_message_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.8 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_star.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.8 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_star_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_wishlist.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_wishlist_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_launcher.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_launcher_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.4 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_message.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "1.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_message_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_star.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.1 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_star_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_wishlist.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_wishlist_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_launcher.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_launcher_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "7.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_message.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_message_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.2 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_star.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.8 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_star_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "6.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_startseite.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_startseite_foreground.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_wishlist.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_wishlist_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_launcher.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "7.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_launcher_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "11.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_message.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.4 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_message_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "8.1 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_star.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "6.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_star_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "10.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_wishlist.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.4 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_wishlist_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "9.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_launcher.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "10.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_launcher_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "16.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_message.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_message_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "11.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_star.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "8.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_star_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "15.2 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_wishlist.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "7.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_wishlist_round.png",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "14.0 kB" }
        },
        {
          "origin": "gradlew",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/unicode.SearchUnicode",
          "acknowledged": false,
          "additionalInfo": {
            "character": "'\\u202a'",
            "filemode": "0100644",
            "filesize": "5.3 kB"
          }
        },
        {
          "origin": "gradle/wrapper/gradle-wrapper.jar",
          "branch": "refs/heads/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "54.3 kB" }
        },
        {
          "origin": "gradlew",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/unicode.SearchUnicode",
          "acknowledged": false,
          "additionalInfo": {
            "character": "'\\u202a'",
            "filemode": "0100644",
            "filesize": "5.3 kB"
          }
        },
        {
          "origin": "app/src/main/ic_message-playstore.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "12.5 kB" }
        },
        {
          "origin": "app/src/main/ic_star-playstore.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "28.0 kB" }
        },
        {
          "origin": "app/src/main/ic_wishlist-playstore.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "20.9 kB" }
        },
        {
          "origin": "app/src/main/res/drawable-v24/micasuca_view_background.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "27.7 kB" }
        },
        {
          "origin": "app/src/main/res/drawable-v24/micasuca_view_home_cut.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "30.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_launcher.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_launcher_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_message.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_message_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.8 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_star.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.8 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_star_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_wishlist.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-hdpi/ic_wishlist_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_launcher.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_launcher_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.4 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_message.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "1.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_message_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_star.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.1 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_star_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_wishlist.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.0 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-mdpi/ic_wishlist_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_launcher.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_launcher_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "7.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_message.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "2.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_message_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.2 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_star.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.8 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_star_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "6.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_startseite.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_startseite_foreground.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_wishlist.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "3.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xhdpi/ic_wishlist_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_launcher.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "7.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_launcher_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "11.9 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_message.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "4.4 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_message_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "8.1 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_star.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "6.3 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_star_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "10.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_wishlist.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.4 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxhdpi/ic_wishlist_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "9.5 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_launcher.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "10.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_launcher_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "16.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_message.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "5.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_message_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "11.7 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_star.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "8.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_star_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "15.2 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_wishlist.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "7.6 kB" }
        },
        {
          "origin": "app/src/main/res/mipmap-xxxhdpi/ic_wishlist_round.png",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "14.0 kB" }
        },
        {
          "origin": "gradle/wrapper/gradle-wrapper.jar",
          "branch": "refs/remotes/origin/master",
          "checkName": "github.com/MrWong99/gitscanner/checks/binaryfile.SearchBinaries",
          "acknowledged": false,
          "additionalInfo": { "filemode": "0100644", "filesize": "54.3 kB" }
        }
      ]
    },
    {
      "date": "2021-12-02T16:16:09.417515+01:00",
      "repository": "git@github.com:MrWong99/IB-documentation.git",
      "error": "",
      "checks": []
    }
  ];

  constructor(
    private messageService: MessageService,
    private searchBinaryService: SearchBinaryService) {
    this.selectedCheckNames.push(this.checknames[0].name);
    
    /*let result = this.fileData.flatMap((data, _, arr) => {      
      data.checks.forEach(e => {
        arr.push(Object.assign({origin: e.origin, branch: e.branch, checkName: e.checkName,
          acknowledged: e.acknowledged, check: null}, data));      
      })
      return arr;
    }, []);    
    console.log(result)*/
  }
 
  ngOnInit() {
  }

  openConfigDialog() {
    this.display = true;
  }

  submit() {   
    let selectedCheckNames: string[] = [];
    this.checknames.forEach((entry) => {
      selectedCheckNames.push(entry.key);
      this.selectedCheckNamesShort = Object.assign([], selectedCheckNames);
    })
    this.getDataSubscription = this.searchBinaryService.getFileData(this.path, selectedCheckNames).subscribe(data => {            
      if (data && data.length > 0) {
        this.fileData = data;
      }
     },
     error => {
      this.messageService.add({key: 'scanError', severity:'error', summary:'Error', detail: error.statusText});
     });
    this.display = false;
  }

  ngOnDestroy(): void {
    this.getDataSubscription?.unsubscribe();
  }
}
